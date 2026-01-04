package csv

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"os"

	"github.com/injoyai/conv"
)

var (
	UTF8    = []byte{0xEF, 0xBB, 0xBF}
	UTF16BE = []byte{0xFE, 0xFF}
	UTF16LE = []byte{0xFF, 0xFE}
)

func Import(filename string) ([][]string, error) {
	result := [][]string(nil)
	err := ImportRange(filename, func(i int, line []string) bool {
		result = append(result, line)
		return true
	})
	return result, err
}

func ImportRange(filename string, fn func(i int, line []string) bool) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	{
		bs, err := buf.Peek(4)
		if err != nil && err != io.EOF {
			if err == io.EOF {
				return nil
			}
			return err
		}
		for _, v := range [][]byte{UTF8, UTF16BE, UTF16LE} {
			if bytes.HasPrefix(bs, v) {
				buf.Discard(len(v))
				break
			}
		}
	}

	r := csv.NewReader(buf)

	for i := 0; ; i++ {
		line, err := r.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if !fn(i, line) {
			return nil
		}
	}
}

func Export[T any](data [][]T) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := buf.Write(UTF8); err != nil {
		return nil, err
	}
	w := csv.NewWriter(buf)
	for _, v := range data {
		if err := w.Write(conv.Strings(v)); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return buf, nil
}
