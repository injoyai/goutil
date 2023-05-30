package in

import (
	"strconv"
)

//返回文件
//@bytes,数据源
//@name,前端显示文件名称
func returnFile(name string, bytes []byte) {
	NewExit(200, bytes).
		SetHeader("Content-Disposition", "attachment; filename="+name).
		SetHeader("Content-Length", strconv.Itoa(len(bytes))).Exit()
}
