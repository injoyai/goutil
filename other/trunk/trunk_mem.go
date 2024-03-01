package trunk

import (
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/logs"
	uuid "github.com/satori/go.uuid"
	"sync"
)

type (
	SubscribeHandler  func(data interface{})
	MiddlewareHandler func(data interface{}) (interface{}, bool)
)

func NewMem() *Mem {
	return &Mem{
		dataQueue:      make(chan interface{}, 1000),
		subscribeQueue: make(chan interface{}, 1000),
	}
}

type Hook struct {
	key       string //用于删除释放通道
	closeFunc func() error
	C         chan interface{}
}

func (this *Hook) Close() error {
	return this.closeFunc()
}

type Mem struct {
	dataQueue      chan interface{}    //数据队列
	middlewareList []MiddlewareHandler //中间件函数

	subscribeQueue chan interface{}   //订阅队列
	subscribeList  []SubscribeHandler //订阅函数
	hookList       []*Hook            //hook函数
	hookLock       sync.Mutex         //hook锁
}

// Publish 发布
func (this *Mem) Publish(data interface{}) {
	this.dataQueue <- data
}

// Subscribe 订阅数据,不会改变数据
func (this *Mem) Subscribe(h SubscribeHandler) {
	this.subscribeList = append(this.subscribeList, h)
}

func (this *Mem) Hook() *Hook {
	this.hookLock.Lock()
	defer this.hookLock.Unlock()
	key := uuid.NewV4().String()
	c := &Hook{
		key: key,
		C:   make(chan interface{}),
		closeFunc: func() error {
			this.hookLock.Lock()
			defer this.hookLock.Unlock()
			for i, v := range this.hookList {
				if v.key == key {
					this.hookList = append(this.hookList[:i], this.hookList[i+1:]...)
					close(v.C)
					break
				}
			}
			return nil
		},
	}
	this.hookList = append(this.hookList, c)
	return c
}

// Middleware 中间件,可以改变数据
func (this *Mem) Middleware(h MiddlewareHandler) {
	this.middlewareList = append(this.middlewareList, h)
}

func (this *Mem) Run() {

	//启一个单独的协程,防止阻塞到主消息队列
	go this.runSubscribe()

loop:
	for data := range this.dataQueue {

		var pass bool
		//执行中间件,能改变原始数据
		for _, h := range this.middlewareList {
			logs.PrintErr(g.Try(func() error {
				data, pass = h(data)
				return nil
			}))
			if !pass {
				//数据不符合,过滤掉
				continue loop
			}
		}

		//尝试加入订阅队列,不影响到主消息队列
		select {
		case this.subscribeQueue <- data:
		default:
		}

		//
		for _, hook := range this.hookList {
			select {
			case hook.C <- data:
			default:
			}
		}
	}
}

func (this *Mem) runSubscribe() {
	for data := range this.subscribeQueue {
		for _, h := range this.subscribeList {
			logs.PrintErr(g.Try(func() error {
				h(data)
				return nil
			}))
		}
	}
}
