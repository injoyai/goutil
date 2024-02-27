package bar

import (
	"io"
	"time"
)

func Demo() error {
	x := New(60)
	x.SetStyle('>')
	x.SetColor(BgCyan)
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			x.Add(1)
		}
	}()
	return x.Run()
}

func Download(url, filename string, proxy ...string) (int64, error) {
	return New(0).DownloadHTTP(url, filename, proxy...)
}

func DownloadHTTP(url, filename string, proxy ...string) (int64, error) {
	return New(0).DownloadHTTP(url, filename, proxy...)
}

func Copy(w io.Writer, r io.Reader, total int64) (int64, error) {
	return New(total).Copy(w, r)
}
