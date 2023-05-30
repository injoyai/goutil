package excel

//import (
//	"fmt"
//	"github.com/tealeg/xlsx"
//	"io"
//	"io/ioutil"
//	"sync"
//)
//
//type Interface interface {
//	Row([]string) (interface{}, error)
//	Insert([]interface{}) ([]string, error)
//}
//
//type Import struct {
//	Interface
//	sync.Mutex
//	read io.Reader
//}
//
//func NewImport(i Interface, r io.Reader) *Import {
//	return &Import{
//		Interface: i,
//		read:      r,
//	}
//}
//
//func (this *Import) Do() ([]string, error) {
//	this.Lock()
//	defer this.Unlock()
//	bs, err := ioutil.ReadAll(this.read)
//	if err != nil {
//		return nil, err
//	}
//	xlFile, err := xlsx.OpenBinary(bs)
//	if err != nil {
//		return nil, err
//	}
//	var msg []string //记录错误信息
//	data := []interface{}{}
//	for i, sheet := range xlFile.Sheets { //页
//		for k, row := range sheet.Rows { //行
//			if k == 0 { //排除第一行
//				continue
//			}
//			cell := []string{}
//			for _, v := range row.Cells {
//				cell = append(cell, v.String())
//			}
//			if len(cell) == 0 {
//				continue
//			}
//			l, err := this.Row(cell)
//			if err != nil {
//				msg = append(msg, fmt.Sprintf("第%d页第%d行:", i, k)+err.Error())
//				continue
//			}
//			if l != nil {
//				data = append(data, l)
//			}
//		}
//	}
//	msg2, err := this.Insert(data)
//	if err != nil {
//		return nil, err
//	}
//	return append(msg, msg2...), nil
//}
