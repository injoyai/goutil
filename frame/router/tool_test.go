package router

import (
	"testing"
)

func Test_cleanPath(t *testing.T) {
	t.Log(cleanPath("/1/2/3/*"))
}
