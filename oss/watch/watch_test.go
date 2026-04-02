package watch

import (
	"testing"

	"github.com/fsnotify/fsnotify"
)

func TestFile(t *testing.T) {
	Watch("./test.txt", func(e Event) {
		t.Log(e)
	})
}

func TestFiles(t *testing.T) {
	err := Watch(
		[]string{
			"./test.txt",
			"./test2.txt",
			"./test/",
			"./test/config.txt",
		},
		func(e fsnotify.Event) {
			t.Log(e)
		},
	)
	t.Log(err)
}
