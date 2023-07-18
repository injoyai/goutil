package task

import (
	"github.com/injoyai/base/chans"
	"github.com/injoyai/goutil/g"
)

type Runner func() error

func NewOnce(retryNum, limit int) *Once {
	return &Once{
		RetryNum: retryNum,
		Limit:    limit,
	}
}

type Once struct {
	RetryNum int                    //重试次数
	Limit    int                    //并发数量
	DoneFunc func(i int, err error) //子项执行完成
}

func (this *Once) SetDoneFunc(f func(i int, err error)) *Once {
	this.DoneFunc = f
	return this
}

func (this *Once) Run(list []Runner) []error {
	errList := []error(nil)
	limit := chans.NewWaitLimit(uint(this.Limit))
	for i, f := range list {
		limit.Add()
		go func(i int, errList []error) {
			defer limit.Done()
			err := g.Retry(f, this.RetryNum)
			if err != nil {
				errList = append(errList, err)
			}
			this.DoneFunc(i, err)
		}(i, errList)
	}
	return errList
}
