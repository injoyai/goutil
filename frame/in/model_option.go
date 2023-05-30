package in

type Func struct {
	Deal func(ok bool, v interface{}, cnt ...int64) interface{} //数据处理
	Succ func(data interface{}, cnt ...int64)
	Fail func(data interface{})
}

// SetSuccFunc 设置成功处理函数,不能为nil
func (this *Func) SetSuccFunc(fn func(interface{}, ...int64)) { this.Succ = fn }

// SetFailFunc 设置失败处理函数,不能为nil
func (this *Func) SetFailFunc(fn func(interface{})) { this.Fail = fn }

// SetDealFunc 设置结果处理函数,不能为nil
func (this *Func) SetDealFunc(fn func(bool, interface{}, ...int64) interface{}) { this.Deal = fn }

type Option struct {
	ExitMark string      //退出标识
	ExitOk   string      //成功退出标示,暂时无效
	CodeSucc interface{} //code,成功
	CodeFail interface{} //code,失败
	Page     string      //页数
	Size     string      //数/页
	SizeDef  int         //默认数量/页
}

// SetExitMark 设置退出标识
func (this *Option) SetExitMark(s string) *Option {
	this.ExitMark = s
	return this
}

// SetExitOk 设置正常退出标识(无效)
func (this *Option) SetExitOk(s string) *Option {
	this.ExitOk = s
	return this
}

// SetCode 设置成功失败标识 "code":succ|fail
func (this *Option) SetCode(succ, fail interface{}) *Option {
	this.CodeSucc = succ
	this.CodeFail = fail
	return this
}

// SetPageSize 设置页数标识,默认页数大小
func (this *Option) SetPageSize(page, size string, sizeDef ...int) *Option {
	this.Page = page
	this.Size = size
	if len(sizeDef) > 0 {
		this.SizeDef = sizeDef[0]
	}
	return this
}

//==============================Default==============================

// defaultFunc 默认函数配置
func defaultFunc() *Func {
	data := &Func{
		Deal: dealWithDefault,
		Succ: succWithDefault,
		Fail: failWithDefault,
	}
	return data
}

// defaultOption 默认选项配置
func defaultOption() *Option {
	return &Option{
		ExitMark: "EXITMARK",
		ExitOk:   "OK",
		CodeSucc: "SUCCESS",
		CodeFail: "FAIL",
		Page:     "index",
		Size:     "size",
		SizeDef:  10,
	}
}

// succWithDefault 默认成功处理函数
func succWithDefault(data interface{}, count ...int64) {
	NewExitJson(200, dealWithDefault(true, data, count...)).Exit()
}

// failWithDefault 默认失败处理函数
func failWithDefault(data interface{}) {
	NewExitJson(200, dealWithDefault(false, data)).Exit()
}

// dealWithDefault 默认结果处理函数
func dealWithDefault(ok bool, v interface{}, cnt ...int64) interface{} {
	m := map[string]interface{}{
		"code": DefaultOption.CodeFail,
		"msg":  "",
		"data": "",
	}
	switch val := v.(type) {
	case error:
		if val != nil {
			m["msg"] = val.Error()
		}
	default:
		if ok {
			m["code"] = DefaultOption.CodeSucc
			m["data"] = v
			if len(cnt) > 0 {
				m["count"] = cnt[0]
			}
		} else if v != nil {
			m["msg"] = v
		}
	}
	return m
}
