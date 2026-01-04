package csv

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"os"

	"github.com/injoyai/conv"
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

		// 1. 匹配 UTF-8 BOM (EF BB BF)
		if bytes.HasPrefix(bs, []byte{0xEF, 0xBB, 0xBF}) {
			buf.Discard(3)
		}

		// 2. 匹配 UTF-16 Big Endian (FE FF)
		if bytes.HasPrefix(bs, []byte{0xFE, 0xFF}) {
			buf.Discard(2)
		}

		// 3. 匹配 UTF-16 Little Endian (FF FE)
		if bytes.HasPrefix(bs, []byte{0xFF, 0xFE}) {
			buf.Discard(2)
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

func Export(data [][]any) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := buf.WriteString("\xEF\xBB\xBF"); err != nil {
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
