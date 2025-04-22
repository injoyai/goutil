package bar

import (
	"io"
	"time"
)

func Demo() error {
	x := New(
		WithTotal(60),
		WithDefaultFormat(func(p Plan) {
			p.SetStyle('>')
			p.SetColor(BgCyan)
		}),
	)
	for {
		time.Sleep(time.Millisecond * 100)
		x.Add(1)
		if x.Flush() {
			return nil
		}
	}
}

func Download(url, filename string, proxy ...string) (int64, error) {
	return New().DownloadHTTP(url, filename, proxy...)
}

func DownloadHTTP(url, filename string, proxy ...string) (int64, error) {
	return New().DownloadHTTP(url, filename, proxy...)
}

func Copy(w io.Writer, r io.Reader, total int64) (int64, error) {
	return New(WithTotal(total)).Copy(w, r)
}
