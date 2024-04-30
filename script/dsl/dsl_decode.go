package dsl

import (
	"encoding/hex"
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
	if len(d.Key) == 0 {
		d.Key = "value"
	}
	for _, v := range opts {
		v(d)
	}
	return d, nil
}

// WithDebug 自定义调试模式
func WithDebug(b ...bool) func(d *Decode) {
	return func(d *Decode) {
		d.Debug = len(b) == 0 || b[0]
	}
}

// WithLog 自定义日志输出
func WithLog(f func(format string, v ...interface{})) func(d *Decode) {
	return func(d *Decode) {
		d.Logger = f
	}
}

// WithKey 自定义key
func WithKey(key string) func(d *Decode) {
	return func(d *Decode) {
		d.Key = key
	}
}

// WithGlobal 自定义全局变量
func WithGlobal(m g.Map) func(d *Decode) {
	return func(d *Decode) {
		d.Global = m
	}
}

// WithScript 自定义脚本
func WithScript(s script.Client) func(d *Decode) {
	return func(d *Decode) {
		d.Script = s
	}
}

func WithDecode(d *Decode) func(d *Decode) {
	return func(d2 *Decode) {
		//这个可以设置成不一样
		d2.Key = d.Key
		//将d的结果赋值到d2
		d2.Global = d.Global
		//将d的脚本赋值到d2,可能设置了脚本的全局配置,也能转移过来
		d2.Script = d.Script
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
	Script  script.Client                         `yaml:"-"`
	Global  g.Map                                 `yaml:"-"`
	index   uint                                  `yaml:"-"`
}

func (this *Decode) newScript(cache g.Map, opts ...func(c script.Client)) script.Client {
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
	if this.Global == nil {
		this.Global = g.Map{this.Key: input}
	}

	//单独的脚本实例
	if this.Script == nil {
		this.Script = this.newScript(this.Global, opts...)
	}

	for _, v := range this.Actions {
		if _, _, err := v.init(this, this.Key).Do(this.Global[this.Key]); err != nil {
			return nil, nil, fmt.Errorf("%v: %v, %v", v.Name, this.Global[this.Key], err)
		}
	}

	delete(this.Global, "_") //同go,匿名变量,不输出
	return this.Global, this.Global[this.Key], nil
}

type Action struct {
	decode    *Decode
	Name      string             `yaml:"name"`   //动作名称
	Key       string             `yaml:"key"`    //全局变量,可选,优先级3
	ParentKey string             `yaml:"-"`      //父级全局变量
	Script    interface{}        `yaml:"script"` //执行脚本,js,优先级1
	Value     interface{}        `yaml:"value"`  //常量,无法通过脚本返回值固增加这个字段,优先级2
	Switch    map[string]*Action `yaml:"switch"` //枚举,默认继承父级的key,优先级5
	Error     string             `yaml:"error"`  //自定义错误信息,优先级4
}

// IsPrivate 是否是私有动作,
// 即设置了key,则当前和子集的数据归这个key所有,否则将赋值到父级
func (this *Action) IsPrivate() bool {
	return len(this.Key) > 0
}

func (this *Action) GetKey() string {
	if len(this.Key) > 0 {
		return this.Key
	}
	return this.ParentKey
}

func (this *Action) GetName() string {
	if len(this.Name) > 0 {
		return this.Name
	}
	return this.GetKey()
}

func (this *Action) init(d *Decode, parentKey string) *Action {
	this.decode = d
	//需要父级的key,而不是顶级的key,固有2个参数
	this.ParentKey = parentKey
	return this
}

// 被一些临时脚本引用,不做赋值处理
func (this *Action) exec(input interface{}) (output interface{}, err error) {
	defer func() {
		if this.decode.Debug {
			fmt.Printf("   - %s: %s, %s: %v, %s: %v\n",
				color.RedString("执行"), this.GetName(),
				color.RedString("输入"), func() interface{} {
					if bs, ok := input.([]byte); ok {
						return "0x" + strings.ToUpper(hex.EncodeToString(bs))
					}
					return input
				}(),
				color.RedString("输出"), output,
			)
		}
	}()

	//todo key和value不是对应关系,获取key不在这里执行

	//加入未设置脚本
	text := conv.String(this.Script)
	if len(text) == 0 {
		switch {
		case len(this.Error) > 0:
			//执行错误信息脚本,例如switch的default
			return nil, this.execError(input)

		case this.Value != nil:
			//设置常量,一般用于switch
			return this.Value, nil

		case this.IsPrivate():
			//只设置了key得话,赋值父级的值到key上
			return input, nil

		default:
			//未设置脚本,未设置错误脚本,未设置常量,未设置key
			return input, nil
		}
	}
	output, err = this.decode.Script.Exec(text, func(i script.Client) {
		for k, v := range this.decode.Global {
			i.Set(k, v)
		}
		//使用上次执行的结果,例如多次执行,能通过value获取上次的结果,相当于闭包的变量
		i.Set(this.decode.Key, input)
	})
	if err != nil {
		return nil, err
	}

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
	if val, ok := output.(bool); ok && len(this.Error) > 0 {
		if !val {
			return output, this.execError(input)
		}
		//设置了错误,不需要返回值
		return nil, nil
	}

	//如果设置了switch,执行的结果赋值进去
	if len(this.Switch) > 0 {
		if sw := this.Switch[conv.String(output)]; sw != nil {
			_, output, err = sw.init(this.decode, this.GetKey()).Do(output)
			return
		}

		//默认,switch的default,如果设置了default,就走default过
		if sw := this.Switch["default"]; sw != nil {
			_, output, err = sw.init(this.decode, this.GetKey()).Do(output)
			return
		}
	}

	return
}

func (this *Action) child() []*Action {
	actions := []*Action(nil)
	switch val := this.Script.(type) {
	case []*Action:
		actions = val

	case nil:
		//处理子集未设置脚本的情况,则是nil
		this.Script = ""
		actions = append(actions, this)

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
				if yaml.Unmarshal(conv.Bytes(v), action) == nil {
					actions = append(actions, action)
				}
			}
		}

	}
	return actions
}

func (this *Action) Do(input interface{}) (key string, output interface{}, err error) {

	this.decode.index++
	fmt.Printf("%s: %s\n - %s: %v\n",
		color.BlueString("步骤%d", this.decode.index), this.GetName(),
		color.RedString("输入"), this.decode.Global,
	)
	defer func() {
		if this.decode.Debug {
			fmt.Printf(" - %s: %v\n\n", color.RedString("输出"), this.decode.Global)
		}
	}()

	//获取key,如果动作设置了key,则使用该动作的key,否则使用父级的key
	//另外key也可以使用脚本生成,例如 "@js getString('filed')"
	key, err = this.execKey(input)
	if err != nil {
		return "", nil, err
	}

	//缓存数据,缓存子集列表的上次执行结果,使子集逻辑通畅
	cache := input
	//默认nil,不受输入的影响,子集全部设置了key,则输出也是nil
	var result interface{}
	for _, v := range this.child() {
		childKey := key
		switch v.Script.(type) {
		case string, nil:
			//不存在子集,直接运行脚本
			result, err = v.init(this.decode, key).exec(cache)
			if err == nil {
				//生成子集的key,子集可能有自己的key
				childKey, err = v.execKey(cache)
			}
		default:
			//存在子集,解析子集,运行子集的结果
			childKey, result, err = v.init(this.decode, key).Do(cache)
		}
		if err != nil {
			return "", nil, err
		}
		if result != nil && v.IsPrivate() {
			//匿名脚本,直接赋值结果,
			//把结果存到缓存中
			//过滤无效的数据,例如只是打印下,则原先的执行结果无法保存
			//过滤掉无效数据后,执行结果能正常保存
			//设置结果到变量,如果设置了key,则赋值到key,否则继承父级的key,即设置到父级
			v.decode.Global[childKey] = result
		}
		if result != nil && !v.IsPrivate() {
			//如果设置了key,则这个结果不输出到父级,即赋值到这个key就是生命周期的终点
			//设置了key的脚本,或者自己脚本,数据归属于这个key
			//类似匿名函数 ,例 key=func()any{xxx}
			cache = result
			output = result
		}
	}

	//过滤无效的数据,例如只是打印下,则原先的执行结果无法保存
	if output != nil {
		this.decode.Global[key] = output
	}

	return
}

func (this *Action) execError(input interface{}) error {
	if strings.HasPrefix(this.Error, "@js") {
		a := &Action{
			decode: this.decode,
			Name:   this.Name + "-Error",
			Key:    this.GetKey(),
			Script: this.Error[3:],
		}
		result, err := a.exec(input)
		if err != nil {
			return err
		}
		return fmt.Errorf("%v", result)
	}
	return errors.New(this.Error)
}

func (this *Action) execKey(input interface{}) (string, error) {
	key := this.GetKey()
	if strings.HasPrefix(key, "@js") {
		a := &Action{
			decode: this.decode,
			Name:   this.Name + "-Key",
			Key:    "_",
			Script: key[3:],
		}
		result, err := a.exec(input)
		if err != nil {
			return "", err
		}
		return conv.String(result), nil
	}
	return key, nil
}
