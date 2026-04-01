package watch

import (
	"testing"

	"github.com/fsnotify/fsnotify"
)

func TestFile(t *testing.T) {
	File("./test.txt", func(e fsnotify.Op) {
		t.Log(e)
	})
}

func TestFiles(t *testing.T) {
	Files(
		[]string{
			"./test.txt",
			"./test2.txt",
			"./test/",
		}, func(e fsnotify.Event) {
			t.Log(e)
		},
	)
}
