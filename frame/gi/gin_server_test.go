package gi

import (
	"github.com/gin-gonic/gin"
	"github.com/injoyai/goutil/frame/in"
	"testing"
)

func TestNew(t *testing.T) {
	s := New(8200)
	s.GET("/test", func(context *gin.Context) {
		in.Fail(nil)
	})
	s.Run()
}
