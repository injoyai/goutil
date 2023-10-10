package zip

import (
	"testing"
)

func Test_unComprZip(t *testing.T) {
	t.Log(Encode("./output", "./output.zip"))
}

func TestDecodeZip(t *testing.T) {
	t.Log(Decode("./output.zip", "./"))
}
