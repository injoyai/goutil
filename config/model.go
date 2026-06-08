package config

const (
	Input  = "input"    //单行文本输入
	Text   = "textarea" //多行文本输入
	Number = "number"   //数字输入框
	Switch = "switch"   //开关
	Select = "select"   //下拉选择
	Color  = "color"    //颜色选择器
	Range  = "range"    //滑块
	Json   = "json"     //JSON 编辑器
	Array  = "array"    //字符串数组
	Dict   = "dict"     //键值对字典
)

type Frame struct {
	Title  string  `json:"title"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Key     string `json:"key"`     //字段名
	Label   string `json:"label"`   //显示名称
	Type    string `json:"type"`    //控件类型
	Default any    `json:"default"` //默认值

	Hint      string    `json:"hint,omitempty"`      //提示
	Options   []Option  `json:"options,omitempty"`   //下拉框专用
	DependsOn DependsOn `json:"dependsOn,omitempty"` //触发条件

	Min    int    `json:"min,omitempty"`    //最小值
	Max    int    `json:"max,omitempty"`    //最大值
	Step   int    `json:"step,omitempty"`   //步长
	Suffix string `json:"suffix,omitempty"` //后缀,例如 "5s" 中的 "s"
}

type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type DependsOn struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}
