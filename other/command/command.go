package command

import (
	"github.com/injoyai/conv"
	"github.com/spf13/cobra"
)

type Command struct {
	cobra.Command                                                      //
	Run           func(cmd *cobra.Command, args []string, flag *Flags) //执行函数,比自带的多个自动解析的flag
	Flag          []*Flag                                              //自动解析flag
	Child         []*Command                                           //设置子项
}

func (this *Command) ParesFlags() *cobra.Command {
	return this.paresFlags()
}

func (this *Command) paresFlags(flags ...*Flag) *cobra.Command {
	for _, v := range this.Flag {
		this.Command.PersistentFlags().StringVarP(&v.Value, v.Name, v.Short, v.Default, v.Memo)
	}
	flags = append(this.Flag, flags...)
	if this.Command.Run == nil && this.Run != nil {
		this.Command.Run = func(cmd *cobra.Command, args []string) {
			if this.Run != nil {
				this.Run(cmd, args, newFlags(flags))
			}
		}
	}
	for _, v := range this.Child {
		this.Command.AddCommand(v.paresFlags(flags...))
	}
	return &this.Command
}

/*



 */

func newFlags(list []*Flag) *Flags {
	f := &Flags{m: make(map[string]*Flag)}
	for _, v := range list {
		f.m[v.Name] = v
	}
	f.Extend = conv.NewExtend(f)
	return f
}

type Flags struct {
	m           map[string]*Flag //保存flag
	conv.Extend                  //继承方法
}

func (this *Flags) Range(fn func(key string, val *Flag) bool) {
	for k, v := range this.m {
		if !fn(k, v) {
			break
		}
	}
}

func (this *Flags) GetVar(key string) *conv.Var {
	val, ok := this.m[key]
	if ok && len(val.Value) > 0 {
		return conv.New(val.Value)
	}
	return conv.Nil()
}

type Flag struct {
	Name    string //flag名称  例如 --name
	Short   string //flag短名称,例如 -n,中的n代替name
	Default string //默认值,如果没有输入的话,使用默认值
	Memo    string //备注信息,提示信息
	Value   string //值,输入的值 --name injoy 中的injoy
}
