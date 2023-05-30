package excel

import (
	"bytes"
	"github.com/injoyai/conv"
	"github.com/tealeg/xlsx"
	"io"
	"io/ioutil"
)

func ToExcel(data [][]interface{}) (*bytes.Buffer, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("sheet1")
	if err != nil {
		return nil, err
	}
	for _, rowValue := range data {
		row := sheet.AddRow()
		for _, cellValue := range rowValue {
			row.AddCell().Value = conv.String(cellValue)
		}
	}
	buf := bytes.NewBuffer(nil)
	if err := file.Write(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func FromExcel(buf io.Reader) (result [][]string, err error) {
	data, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, err
	}
	file, err := xlsx.OpenBinary(data)
	if err != nil {
		return nil, err
	}
	sheet := file.Sheets[0]
	for _, row := range sheet.Rows {
		slice := make([]string, 0)
		for _, cell := range row.Cells {
			slice = append(slice, cell.Value)
		}
		result = append(result, slice)
	}
	return
}
