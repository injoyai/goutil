package zip

import (
	"testing"
)

func Test_unComprZip(t *testing.T) {
	t.Log(Encode("./test/", "./output.zip"))
}

func TestDecodeZip(t *testing.T) {
	t.Log(Decode("./output.zip", "./test/"))
}

func Test_unComprZip2(t *testing.T) {
	t.Log(Encode("./chrome", "./output.zip"))
}

func TestDecodeZip2(t *testing.T) {
	t.Log(Decode("./chrome.zip", "."))
}
