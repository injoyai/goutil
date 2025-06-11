package task

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"sync"
	"time"
)

// New 新建计时器(任务调度),最小周期秒
func New[K comparable]() *Cron[K] {
	return &Cron[K]{
		c: cron.New(cron.WithSeconds()),
		m: make(map[K]*Task[K]),
	}
}

// Cron 定时器(任务调度),任务起一个协程
type Cron[K comparable] struct {
	c  *cron.Cron
	m  map[K]*Task[K]
	mu sync.RWMutex
}

func (this *Cron[K]) Start() *Cron[K] {
	this.c.Start()
	return this
}

func (this *Cron[K]) Run() {
	this.c.Run()
}

func (this *Cron[K]) Stop() *Cron[K] {
	this.c.Stop()
	return this
}

// GetTaskAll 读取全部任务
func (this *Cron[K]) GetTaskAll() []*Task[K] {
	m := make(map[cron.EntryID]*Task[K])
	this.mu.RLock()
	for _, v := range this.m {
		m[v.ID] = v
	}
	this.mu.RUnlock()

	taskList := []*Task[K](nil)
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
func (this *Cron[K]) GetTask(key K) *Task[K] {
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
func (this *Cron[K]) SetTask(key K, spec string, handler func()) error {
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
func (this *Cron[K]) DelTask(key K) {
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
type Task[K comparable] struct {
	Key        K      //任务唯一标识
	Spec       string //定时规则
	cron.Entry        //任务
}

func newTask[K comparable](key K, spec string, e cron.Entry) *Task[K] {
	return &Task[K]{
		Key:   key,
		Spec:  spec,
		Entry: e,
	}
}

func (this *Task[K]) SetEntry(e cron.Entry) *Task[K] {
	this.Entry = e
	return this
}

func (this *Task[K]) String() string {
	return fmt.Sprintf("名称(%v),生效(%v),规则(%s),上次执行时间(%v),下次执行时间(%v)",
		this.Key, this.Valid(), this.Spec, this.timeStr(this.Prev), this.timeStr(this.Next))
}

func (this *Task[K]) timeStr(t time.Time) string {
	if t.IsZero() {
		return "无"
	}
	return t.String()
}
