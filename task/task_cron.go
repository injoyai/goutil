package task

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"sync"
	"time"
)

type Handler func() error

// New 新建计时器(任务调度),最小周期秒
func New() *Cron {
	return &Cron{
		c: cron.New(cron.WithSeconds()),
		m: make(map[string]*Task),
	}
}

// Cron 定时器(任务调度),任务起一个协程
type Cron struct {
	c  *cron.Cron
	m  map[string]*Task
	mu sync.RWMutex
}

func (this *Cron) Start() *Cron {
	this.c.Start()
	return this
}

func (this *Cron) Stop() *Cron {
	this.c.Stop()
	return this
}

// GetTaskAll 读取全部任务
func (this *Cron) GetTaskAll() []*Task {
	m := make(map[cron.EntryID]*Task)
	this.mu.RLock()
	for _, v := range this.m {
		m[v.ID] = v
	}
	this.mu.RUnlock()

	taskList := []*Task(nil)
	for _, v := range this.c.Entries() {
		task, ok := m[v.ID]
		if !ok {
			//删除不在档的任务
			this.c.Remove(v.ID)
			continue
		}
		taskList = append(taskList, task.SetEntry(v))
	}
	return taskList
}

// GetTask 读取任务
func (this *Cron) GetTask(key string) *Task {
	this.mu.RLock()
	task, ok := this.m[key]
	this.mu.RUnlock()
	if !ok {
		return nil
	}
	en := this.c.Entry(task.ID)
	if en.ID == 0 {
		this.mu.Lock()
		delete(this.m, key)
		this.mu.Unlock()
		return nil
	}
	return task.SetEntry(en)
}

// SetTask 设置任务
func (this *Cron) SetTask(key, spec string, handler func()) error {
	this.mu.RLock()
	task, ok := this.m[key]
	this.mu.RUnlock()
	if ok {
		//存在相同任务,则移除
		this.c.Remove(task.ID)
		this.mu.Lock()
		delete(this.m, key)
		this.mu.Unlock()
	}
	id, err := this.c.AddFunc(spec, handler)
	if err != nil {
		return err
	}
	this.mu.Lock()
	this.m[key] = newTask(key, spec, cron.Entry{ID: id})
	this.mu.Unlock()
	return nil
}

// DelTask 删除任务
func (this *Cron) DelTask(key string) {
	this.mu.RLock()
	task, ok := this.m[key]
	this.mu.RUnlock()
	if ok {
		this.c.Remove(task.ID)
		this.mu.Lock()
		delete(this.m, key)
		this.mu.Unlock()
	}
}

// Task 任务
type Task struct {
	Key        string //任务唯一标识
	Spec       string //定时规则
	cron.Entry        //任务
}

func newTask(key, spec string, e cron.Entry) *Task {
	return &Task{
		Key:   key,
		Spec:  spec,
		Entry: e,
	}
}

func (this *Task) SetEntry(e cron.Entry) *Task {
	this.Entry = e
	return this
}

func (this *Task) String() string {
	return fmt.Sprintf("名称(%s),生效(%v),规则(%s),上次执行时间(%v),下次执行时间(%v)",
		this.Key, this.Valid(), this.Spec, this.timeStr(this.Prev), this.timeStr(this.Next))
}

func (this *Task) timeStr(t time.Time) string {
	if t.IsZero() {
		return "无"
	}
	return t.String()
}
