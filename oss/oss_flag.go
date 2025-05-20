package oss

import (
	"flag"
	"github.com/injoyai/conv"
)

func NewFlagConv(flags ...*Flag) conv.Extend {
	return conv.NewExtend(NewFlag(flags...))
}

func NewFlag(flags ...*Flag) *Flags {
	for _, v := range flags {
		flag.String(v.Name, conv.String(v.Default), v.Usage)
	}
	flag.Parse()
	return &Flags{}
}

type Flags struct{}

func (this *Flags) GetVar(key string) *conv.Var {
	f := flag.Lookup(key)
	if f == nil || f.Value.String() == "" {
		return conv.Nil()
	}
	return conv.New(f.Value.String())
}

type Flag struct {
	Name    string //名称
	Default any    //默认值
	Usage   string //使用说明
}
