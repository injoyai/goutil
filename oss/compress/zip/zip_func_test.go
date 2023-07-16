package zip

import (
	"testing"
)

func Test_unComprZip(t *testing.T) {

}

func TestDecodeZip(t *testing.T) {
	t.Log(Decode("./test/test.zip", "./test/test2"))
}
