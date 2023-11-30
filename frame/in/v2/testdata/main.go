package main

import (
	"github.com/gin-gonic/gin"
	in "github.com/injoyai/goutil/frame/in/v2"
)

func main() {
	s := gin.Default()
	in.DefaultClient.CodeSucc = 888
	in.DefaultClient.InitGin(s)
	in.DefaultClient.CodeFail = "FAIL"

	s.GET("/ping/1", func(context *gin.Context) {
		panic(7)
		in.Fail(1)
	})
	s.Run(":8080")
}
