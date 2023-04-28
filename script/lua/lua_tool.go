package lua

import (
	"fmt"
	"strconv"
	"strings"
)

func dealErr(err error) error {
	if err != nil {
		errMsg := err.Error()
		errMsg = strings.Split(errMsg, "\n")[0]

		if list := strings.Split(errMsg, ":"); len(list) == 3 {
			newErr := &Err{Msg: strings.TrimSpace(list[2])}
			newErr.Line, _ = strconv.Atoi(list[1])
			switch newErr.Msg {
			case "attempt to call a non-function object":
				newErr.Msg = "使用不存在(未声明)的函数"
			default:
				newErr.Msg = strings.Replace(newErr.Msg, "cannot perform ", "无法执行", 1)
				newErr.Msg = strings.Replace(newErr.Msg, "mul", "乘法", 1)
				newErr.Msg = strings.Replace(newErr.Msg, "div", "除法", 1)
				newErr.Msg = strings.Replace(newErr.Msg, "add", "加法", 1)
				newErr.Msg = strings.Replace(newErr.Msg, "sub", "减法", 1)
				newErr.Msg = strings.Replace(newErr.Msg, " operation between ", "操作", 1)
				newErr.Msg = strings.Replace(newErr.Msg, " and ", "和", 1)
			}
			return newErr
		}
	}
	return err
}

type Err struct {
	Line int
	Msg  string
}

func (this *Err) Error() string {
	return fmt.Sprintf("行数:%d :%s", this.Line, this.Msg)
}
