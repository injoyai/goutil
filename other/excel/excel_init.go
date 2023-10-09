package excel

import (
	"bytes"
	"github.com/injoyai/conv"
	"github.com/tealeg/xlsx"
	"io"
	"io/ioutil"
)

func ToExcel(sheets map[string][][]interface{}) (*bytes.Buffer, error) {
	file := xlsx.NewFile()
	for sheetName, data := range sheets {
		sheet, err := file.AddSheet(sheetName)
		if err != nil {
			return nil, err
		}
		for _, rowValue := range data {
			row := sheet.AddRow()
			for _, cellValue := range rowValue {
				row.AddCell().Value = conv.String(cellValue)
			}
		}
	}
	buf := bytes.NewBuffer(nil)
	if err := file.Write(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func FromExcel(buf io.Reader) (result map[string][][]string, err error) {
	data, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, err
	}
	file, err := xlsx.OpenBinary(data)
	if err != nil {
		return nil, err
	}
	for _, sheet := range file.Sheets {
		rows := make([][]string, 0, len(sheet.Rows))
		for _, row := range sheet.Rows {
			cells := make([]string, 0, len(row.Cells))
			for _, cell := range row.Cells {
				cells = append(cells, cell.Value)
			}
			rows = append(rows, cells)
		}
		result[sheet.Name] = rows
	}
	return
}
