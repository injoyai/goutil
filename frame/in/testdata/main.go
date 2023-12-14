package main

import (
	"github.com/gin-gonic/gin"
	"github.com/injoyai/goutil/frame/in"
	"github.com/injoyai/goutil/g"
	"io"
)

func main() {
	s := gin.Default()
	in.DefaultClient.InitGin(s)
	//in.DefaultClient.SetCode("SUCCESS", "FAIL")

	s.GET("/1", func(context *gin.Context) {
		panic(7)
		in.Fail(1)
	})
	s.GET("/2", func(context *gin.Context) {
		in.Fail(1)
	})
	s.GET("/3", func(context *gin.Context) {
		in.Succ(g.Map{
			"a": 1, "b": true, "c": 3.2, "d": "s",
		})
	})
	s.GET("/4", func(context *gin.Context) {
		in.Json415(415)
	})
	s.Any("/get", func(context *gin.Context) {

		x := in.GetString(context, "x")

		x2 := in.GetBodyMap(context).GetString("x")
		bs, _ := io.ReadAll(context.Request.Body)
		in.Succ(g.Map{
			"x":    x,
			"x2":   x2,
			"body": string(bs),
		})
	})
	s.Any("/proxy", func(context *gin.Context) {
		in.Proxy(context.Writer, context.Request, "https://www.baidu.com")
	})
	s.Any("/file", func(context *gin.Context) {
		in.FileLocal("test.txt", "./main.go")
	})
	s.Any("/text", func(context *gin.Context) {
		in.Text(200, "666")
	})
	s.Run(":8080")
}
