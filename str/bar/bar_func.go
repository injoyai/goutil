package bar

import (
	"github.com/fatih/color"
	"io"
	"time"
)

func Demo() <-chan struct{} {
	x := New(60)
	x.SetStyle('>')
	x.SetColor(color.BgCyan)
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			x.Add(1)
		}
	}()
	return x.Run()
}

func Download(url, filename string, proxy ...string) error {
	return New(0).DownloadHTTP(url, filename, proxy...)
}

func DownloadHTTP(url, filename string, proxy ...string) error {
	return New(0).DownloadHTTP(url, filename, proxy...)
}

func Copy(w io.Writer, r io.Reader, total int64) error {
	return New(total).Copy(w, r)
}
