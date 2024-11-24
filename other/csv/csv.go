package csv

import (
	"bufio"
	"encoding/csv"
	"os"
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
			return err
		}
		if !fn(line) {
			return nil
		}
	}
}

func Export(data [][]interface{}) (*bufio.Reader, error) {

	return nil, nil
}
