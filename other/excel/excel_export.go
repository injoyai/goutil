package excel

import (
	"bytes"
	"github.com/injoyai/conv"
	"github.com/tealeg/xlsx"
)

type Export map[string][][]any //数据源,页,行,列

// Set 添加数据
// @data,数据源
// @sheetName,分页名称,可选(默认sheet1)
func (this *Export) Set(data [][]any, sheetName ...string) *Export {
	if this == nil {
		*this = make(map[string][][]any)
	}
	name := "Sheet1"
	if len(sheetName) != 0 && len(sheetName[0]) != 0 {
		name = sheetName[0]
	}
	(*this)[name] = data
	return this
}

// Add 添加数据
// @data,数据源
// @sheetName,分页名称,可选(默认sheet1)
func (this *Export) Add(data []any, sheetName ...string) *Export {
	if this == nil {
		*this = make(map[string][][]any)
	}
	name := "Sheet1"
	if len(sheetName) != 0 && len(sheetName[0]) != 0 {
		name = sheetName[0]
	}
	(*this)[name] = append((*this)[name], data)
	return this
}

func (this *Export) Buffer() (*bytes.Buffer, error) {
	file := xlsx.NewFile()
	for sheetName, data := range *this {
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
