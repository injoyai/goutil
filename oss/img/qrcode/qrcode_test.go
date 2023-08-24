package qrcode

import (
	"testing"
)

func TestEncodeLocal(t *testing.T) {
	EncodeLocal("./test.png", "http://aiot.qianlangtech.com/docs")
}
