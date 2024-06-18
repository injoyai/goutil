package oss

import (
	"github.com/injoyai/conv"
	"os"
)

func NewEnvConv() conv.Extend {
	return conv.NewExtend(NewEnv())
}

func NewEnv() *Env { return &Env{} }

type Env struct{}

func (this *Env) GetVar(key string) *conv.Var {
	if v, ok := os.LookupEnv(key); ok {
		return conv.New(v)
	}
	return conv.Nil()
}
