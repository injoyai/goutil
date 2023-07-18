package task

import (
	"github.com/injoyai/base/chans"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/g"
)

func NewRange(limit int, retry ...int) *Range {
	return &Range{
		Retry: conv.GetDefaultInt(3, retry...),
		Limit: limit,
	}
}

type Range struct {
	Retry    int                    //重试次数
	Limit    int                    //并发数量
	DoneFunc func(i int, err error) //子项执行完成
}

func (this *Range) SetRetry(retry int) *Range {
	this.Retry = retry
	return this
}

func (this *Range) SetLimit(limit int) *Range {
	this.Limit = limit
	return this
}

func (this *Range) SetDoneFunc(f func(i int, err error)) *Range {
	this.DoneFunc = f
	return this
}

func (this *Range) Run(list []Handler) []error {
	errList := []error(nil)
	limit := chans.NewWaitLimit(uint(this.Limit))
	for i, f := range list {
		limit.Add()
		go func(i int, errList []error) {
			defer limit.Done()
			err := g.Retry(f, this.Retry)
			if err != nil {
				errList = append(errList, err)
			}
			this.DoneFunc(i, err)
		}(i, errList)
	}
	return errList
}
