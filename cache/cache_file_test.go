package cache

import (
	"testing"
)

func TestFile_getPath(t *testing.T) {
	t.Log(newFile("name", "tag").filename())
	DefaultDir = "./a/b"
	t.Log(newFile("name", "tag").filename())
	t.Log(newFile("./a/b/name", "tag").filename())
}
