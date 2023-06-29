package bar

import (
	"bufio"
	"github.com/fatih/color"
	"github.com/injoyai/conv"
	"io"
	"net/http"
	"os"
	"time"
)

func Demo() <-chan struct{} {
	x := New(60)
	x.SetPrefix("进度: ")
	x.SetSuffix("$")
	x.SetStyle('#')
	x.SetColor(color.FgBlue)
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			x.Add(1)
		}
	}()
	return x.Run()
}

func Download(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return Copy(f, resp.Body, conv.Int64(resp.Header.Get("Content-Length")))
}

func Copy(w io.Writer, r io.Reader, total int64) error {
	buff := bufio.NewReader(r)
	b := New(total)
	b.SetTotal(total)
	go b.Run()
	defer b.Done()
	for {
		buf := make([]byte, 1<<20)
		n, err := buff.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		b.Add(int64(n))
		if _, err := w.Write(buf[:n]); err != nil {
			return err
		}
		if err == io.EOF {
			return nil
		}
	}
}
