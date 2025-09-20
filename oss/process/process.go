package process

import (
	"bufio"
	"context"
	"fmt"
	"github.com/injoyai/base/str"
	"github.com/injoyai/goutil/g"
	"io"
	"os/exec"
	"sync"
)

/*
Process 一个可控制的子进程
*/
type Process interface {
	Close() error
	Running() bool
	Run(ctx ...context.Context) error
	Rerun(ctx ...context.Context) error
	SetStdout(w io.Writer)
	SetStderr(w io.Writer)
}

func New(command string, args ...string) Process {
	return &process{
		args: args,
		cmd:  exec.Command(command, args...),
	}
}

// OutputHandler 是输出回调函数类型
// prefix: "STDOUT"/"STDERR"
// line: 子进程输出的一行
type OutputHandler func(prefix, line string)

// 管理一个可控制的子进程
type process struct {
	cmd     *exec.Cmd
	args    []string
	running bool
	lock    sync.Mutex
	stdout  io.Writer
	stderr  io.Writer
}

func (p *process) SetStdout(w io.Writer) {
	p.stdout = w
}

func (p *process) SetStderr(w io.Writer) {
	p.stderr = w
}

func (p *process) Running() bool {
	return p.running
}

func (p *process) Run(ctx ...context.Context) error {
	if err := p.Start(ctx...); err != nil {
		return err
	}
	return p.cmd.Wait()
}

// Start 启动子进程
func (p *process) Start(ctx ...context.Context) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if p.running {
		return fmt.Errorf("process already running")
	}

	p.cmd = exec.CommandContext(g.Ctx(ctx...), p.cmd.Path, p.args...)
	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = p.cmd.Start()
	if err != nil {
		return err
	}
	p.running = true

	// 输出处理
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Bytes()
			if p.stdout != nil {
				line, _ = str.GbkToUtf8(line)
				p.stdout.Write(line)
			}
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Bytes()
			if p.stderr != nil {
				line, _ = str.GbkToUtf8(line)
				p.stderr.Write(line)
			}
		}
	}()

	return nil
}

// Close 强制关闭子进程
func (p *process) Close() error {
	if !p.running || p.cmd.Process == nil {
		return nil
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	err := p.cmd.Process.Kill()
	if err != nil {
		return err
	}
	p.cmd.Wait()
	p.running = false
	return nil
}

// Rerun 重启子进程
func (p *process) Rerun(ctx ...context.Context) error {
	if p.running {
		p.Close()
	}
	return p.Run(ctx...)
}
