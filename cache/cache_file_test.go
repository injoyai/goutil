package cache

import (
	"testing"
)

func TestFile_getPath(t *testing.T) {
	t.Log(newFile("name", "tag").Filename())
	DefaultDir = "./a/b"
	t.Log(newFile("name", "tag").Filename())
	t.Log(newFile("./a/b/name", "tag").Filename())
}
