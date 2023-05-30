package in

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

//========================== 文本 =======================

//Return 返回字符串
func Return(code int, data string) { NewExit(code, []byte(data)).Exit() }

//Return200 返回字符串,状态码200
func Return200(data string) { Return(200, data) }

// Redirect 重定向
func Redirect(addr string) { NewExit(302, nil).SetHeader("Location", addr).Exit() }

//========================== 文件 =======================

//ReturnFilePath 返回本地文件
func ReturnFilePath(name string, path string) {
	ReturnFileLocal(name, path)
}

//ReturnFileLocal 返回本地文件
func ReturnFileLocal(name string, path string) {
	file, err := os.Open(path)
	CheckErr(err)
	ReturnFileReader(name, file)
}

//ReturnFileReader 返回文件
func ReturnFileReader(name string, read io.Reader) {
	bs, err := ioutil.ReadAll(read)
	CheckErr(err)
	returnFile(name, bs)
}

//ReturnFileBytes 返回文件
func ReturnFileBytes(name string, bytes []byte) {
	returnFile(name, bytes)
}

//========================== json =======================

//Json 返回json,根基方法
func Json(code int, ok bool, data interface{}, count ...int64) {
	NewExitJson(code, DefaultFunc.Deal(ok, data, count...)).Exit()
}

//Json200 返回json,状态码200
func Json200(data interface{}, count ...int64) { Json(200, true, data, count...) }

//Json400 返回json,状态码400
func Json400(data interface{}) { Json(400, false, data) }

//Json401 返回json,状态码401
func Json401() { Json(401, false, "验证失败") }

//Json403 返回json,状态码403
func Json403() { Json(403, false, "没有权限") }

//Json415 返回json,状态码415
func Json415(data interface{}) { Json(415, false, data) }

//Json500 返回json,状态码500
func Json500(data interface{}) { Json(500, false, data) }

//Succ 成功,可配置
func Succ(data interface{}, count ...int64) { DefaultFunc.Succ(data, count...) }

//Fail 失败,可配置
func Fail(data interface{}) { DefaultFunc.Fail(data) }

//Fail200 失败,状态码200
func Fail200(data interface{}) { Json(200, false, data) }

// Errf 退出格式化错误信息
func Errf(format string, args ...interface{}) {
	Err(fmt.Sprintf(format, args...))
}

//Err 退出,并校验错误
func Err(data interface{}) {
	if data == nil {
		Succ(data)
	} else {
		Fail(data)
	}
}

//CheckErr 检测错误(遇到错误结束)
func CheckErr(err error, msg ...string) {
	if err != nil {
		if len(msg) > 0 {
			Err(msg)
		}
		Err(err)
	}
}
