package excel

import (
	"bytes"
	"encoding/csv"
	"github.com/injoyai/conv"
	"github.com/tealeg/xlsx"
	"io"
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

func From(i interface{}) (result map[string][][]string, err error) {
	return FromBytes(conv.Bytes(i))
}

func FromReader(r io.Reader) (result map[string][][]string, err error) {
	bs, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return FromBytes(bs)
}

func FromBytes(bs []byte) (result map[string][][]string, err error) {
	file, err := xlsx.OpenBinary(bs)
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

/*



 */

func FromCsv(i interface{}) (result [][]string, err error) {
	return FromCsvReader(bytes.NewReader(conv.Bytes(i)))
}

func FromCsvReader(r io.Reader) (result [][]string, err error) {
	x := csv.NewReader(r)
	x.FieldsPerRecord = -1
	result, err = x.ReadAll()
	return
}

func FromCsvBytes(bs []byte) (result [][]string, err error) {
	return FromCsvReader(bytes.NewReader(bs))
}

func ToCsv(data [][]interface{}) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(buf)
	for _, rows := range data {
		if err := w.Write(conv.Strings(rows)); err != nil {
			return nil, err
		}
	}
	w.Flush()
	return buf, nil
}
