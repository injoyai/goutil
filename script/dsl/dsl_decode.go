package dsl

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/goutil/script"
	"github.com/injoyai/goutil/script/js"
	"gopkg.in/yaml.v3"
	"strconv"
	"strings"
)

func NewDecode(resource interface{}, opts ...func(*Decode)) (*Decode, error) {
	d := new(Decode)
	if err := yaml.Unmarshal(conv.Bytes(resource), d); err != nil {
		return nil, err
	}
	for _, v := range opts {
		v(d)
	}
	if len(d.Key) == 0 {
		d.Key = "value"
	}
	return d, nil
}

func WithDebug(b ...bool) func(d *Decode) {
	return func(d *Decode) {
		d.Debug = len(b) == 0 || b[0]
	}
}

func WithLog(f func(format string, v ...interface{})) func(d *Decode) {
	return func(d *Decode) {
		d.Logger = f
	}
}

func WithKey(key string) func(d *Decode) {
	return func(d *Decode) {
		d.Key = key
	}
}

/*
Decode 协议解析
解析速度受到步骤的影响
通过脚本实现,增加了灵活性,降低了执行速度(约0.4ms/次),正常使用够了
耗时:

	DLT645(电量解析)  (3.7,4.0,3.7,4.4,4.1,4.1,3.7)s/万次
*/
type Decode struct {
	Name    string                                `yaml:"name"`    //名称
	Actions []*Action                             `yaml:"actions"` //动作
	Debug   bool                                  `yaml:"debug"`   //调试模式
	Logger  func(format string, v ...interface{}) `yaml:"-"`       //日志
	Key     string                                `yaml:"key"`     //脚本取值的key,默认value,例如cut(value,0,2),直接操作变量
}

func (this *Decode) script(cache g.Map, opts ...func(c script.Client)) script.Client {
	s := js.New(script.WithBaseFunc)
	s.Set("get", func(args *script.Args) interface{} {
		return cache[args.Get(1).String()]
	})
	s.Set("getString", func(args *script.Args) interface{} {
		return conv.String(cache[args.Get(1).String()])
	})
	s.Set("getBool", func(args *script.Args) interface{} {
		return conv.Bool(cache[args.Get(1).String()])
	})
	s.Set("getFloat", func(args *script.Args) interface{} {
		return conv.Float64(cache[args.Get(1).String()])
	})
	s.Set("getInt", func(args *script.Args) interface{} {
		return conv.Int64(cache[args.Get(1).String()])
	})
	s.Set("set", func(args *script.Args) {
		cache[args.Get(1).String()] = args.Get(1)
	})
	s.Set("del", func(args *script.Args) {
		for _, v := range args.Args {
			delete(cache, v.String())
		}
	})
	for _, v := range opts {
		v(s)
	}
	return s
}

func (this *Decode) Do(input interface{}, opts ...func(c script.Client)) (g.Map, interface{}, error) {

	//缓存输入输出结果数据
	//设置父级的key为"value",如果动作的key未设置,则会赋值到父级数据上,是所预期的结果
	if len(this.Key) == 0 {
		this.Key = "value"
	}
	cache := g.Map{this.Key: input}

	//单独的脚本实例
	script := this.script(cache, opts...)

	a := &Action{
		decode: this,
		Name:   this.Name,
		Key:    this.Key,
		Script: this.Actions,
	}
	result, err := a.Do(script, cache, input)
	if err != nil {
		return nil, nil, fmt.Errorf("%v: %v, %v", this.Name, cache[this.Key], err)
	}

	//for _, v := range this.Actions {
	//	if err := v.init(this, this.Key).Do(script, cache, cache[this.Key]); err != nil {
	//		return nil, nil, fmt.Errorf("%v: %v, %v", v.Name, cache[this.Key], err)
	//	}
	//}

	return cache, result, nil
}

type Action struct {
	decode *Decode
	Name   string             `yaml:"name"`   //动作名称
	Key    string             `yaml:"key"`    //全局变量,可选,优先级3
	Script interface{}        `yaml:"script"` //执行脚本,js,优先级1
	Value  interface{}        `yaml:"value"`  //常量,无法通过脚本返回值固增加这个字段,优先级2
	Switch map[string]*Action `yaml:"switch"` //枚举,默认继承父级的key,优先级5
	Error  string             `yaml:"error"`  //自定义错误信息,优先级4
}

func (this *Action) init(d *Decode, parentKey string) *Action {
	this.decode = d
	if len(this.Key) == 0 {
		this.Key = parentKey
	}
	return this
}

func (this *Action) exec(s script.Client, global g.Map, scripts []string, parentValue interface{}) (interface{}, error) {
	//加入未设置脚本
	if len(scripts) == 0 {
		switch {
		case len(this.Error) > 0:
			//执行错误信息脚本,例如switch的default
			return nil, this.execError(s, global, parentValue)

		case this.Value != nil:
			//设置常量,一般用于switch
			return this.Value, nil

		case len(this.Key) > 0:
			//只设置了key得话,赋值父级的值到key上
			return parentValue, nil

		default:
			//未设置脚本,未设置错误脚本,未设置常量,未设置key
			return parentValue, nil
		}
	}
	var result interface{}
	for _, text := range scripts {
		output, err := s.Exec(text, func(i script.Client) {
			for k, v := range global {
				i.Set(k, v)
			}
			//使用上次执行的结果,例如多次执行,能通过value获取上次的结果,相当于闭包的变量
			i.Set(this.decode.Key, parentValue)
		})
		if err != nil {
			return nil, err
		}
		if output != nil {
			//过滤无效的数据,例如只是打印下,则原先的执行结果无法保存
			//过滤掉无效数据后,执行结果能正常保存
			result = output
		}
	}
	return result, nil
}

func (this *Action) Do(script script.Client, global g.Map, parentValue interface{}) (result interface{}, err error) {

	defer func() {
		if this.decode.Debug {
			fmt.Printf("%s: %s, %s: %v, %s: %v, %s: %v\n",
				color.RedString("动作"), this.Name,
				color.RedString(this.decode.Key), parentValue,
				color.RedString("输出"), global,
				color.RedString("结果"), conv.New(err).String("成功"))
		}
	}()

	//获取key,如果动作设置了key,则使用该动作的key,否则使用父级的key
	//另外key也可以使用脚本生成,例如 "@js getString('filed')"
	key, err := this.execKey(script, global, parentValue)
	if err != nil {
		return nil, err
	}

	//执行脚本
	actions := []*Action(nil)
	switch val := this.Script.(type) {
	case []*Action:
		actions = val
	case string:
		actions = append(actions, this)
	case []interface{}:
		for i, v := range val {
			switch s := v.(type) {
			case string:
				actions = append(actions, &Action{
					decode: this.decode,
					Name:   this.Name + "-" + strconv.Itoa(i+1),
					Script: s,
				})
			case map[string]interface{}:
				action := new(Action)
				if yaml.Unmarshal(conv.Bytes(v), v) == nil {
					actions = append(actions, action)
				}
			}
		}
	}

	for _, v := range actions {
		result, err = v.init(this.decode, key).Do(script, global, parentValue)
		if err != nil {
			return nil, err
		}
	}

	//result, err := this.exec(script, global, conv.Strings(this.Script), parentValue)
	//if err != nil {
	//	return err
	//}

	/*
		返回bool类型,并值是false,并设置了错误信息,则返回错误
		简化前步骤,正常需要设置临时的key,然后执行错误脚本判断,合并这2个步骤
		- name: 判断xxx
		  key: succ
		  script: 1!=2

		- error: if (succ) {return "错误信息"}

		简化后
		- name: 判断xxx
		  script: 1!=2
		  error: 错误信息
	*/
	if val, ok := result.(bool); ok && len(this.Error) > 0 {
		if !val {
			return result, this.execError(script, global, parentValue)
		}
		return result, nil
	}

	//如果设置了switch,执行的结果赋值进去
	if len(this.Switch) > 0 {
		if sw := this.Switch[conv.String(result)]; sw != nil {
			return sw.init(this.decode, this.Key).Do(script, global, result)
		}

		//默认,switch的default,如果设置了default,就走default过
		if sw := this.Switch["default"]; sw != nil {
			return sw.init(this.decode, this.Key).Do(script, global, result)
		}
	}

	// 把结果存到缓存中
	if result != nil {
		global[key] = result
	}

	return
}

func (this *Action) execError(script script.Client, global g.Map, parentValue interface{}) error {
	if strings.HasPrefix(this.Error, "@js") {
		result, err := this.exec(script, global, []string{this.Error[3:]}, parentValue)
		if err != nil {
			return err
		}
		return fmt.Errorf("%v", result)
	}
	return errors.New(this.Error)
}

func (this *Action) execKey(script script.Client, global g.Map, parentValue interface{}) (string, error) {
	if strings.HasPrefix(this.Key, "@js") {
		result, err := this.exec(script, global, []string{this.Key[3:]}, parentValue)
		if err != nil {
			return "", err
		}
		return conv.String(result), nil
	}
	return this.Key, nil
}
