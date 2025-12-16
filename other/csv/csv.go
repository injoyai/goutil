package csv

import (
	"bytes"
	"encoding/csv"
	"io"
	"os"

	"github.com/injoyai/conv"
)

func Import(filename string, fn func(line []string) bool) ([][]string, error) {
	result := [][]string(nil)
	err := ImportRange(filename, func(line []string) bool {
		result = append(result, line)
		return true
	})
	return result, err
}

func ImportRange(filename string, fn func(line []string) bool) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	r := csv.NewReader(f)
	for {
		line, err := r.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if !fn(line) {
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
