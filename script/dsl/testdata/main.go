package main

import (
	"encoding/hex"
	"fmt"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/goutil/script"
	"github.com/injoyai/goutil/script/dsl"
	"github.com/injoyai/logs"
	"os"
	"path/filepath"
)

func main() {

	dlt645()
	modbusRtu()

}

func modbusRtu() {
	do("./script/dsl/testdata/modbus.yaml", []byte{0x01, 0x03, 0x04, 0x00, 0x19, 0x00, 0x30, 0xe0, 0x2b})
}

func dlt645() {
	do("./script/dsl/testdata/dlt645.yaml",
		g.Map{
			"code": 200,
			"data": g.Map{
				"version": "v1.0",
				"hex":     "68AAAAAAAAAAAA68910833333333343333337E16",
			},
		},
	)
}

func do(filename string, input any) {
	fmt.Println("\n\n==============================" + filepath.Base(filename) + "==============================")
	bs, err := os.ReadFile(filename)
	logs.PanicErr(err)
	d, err := dsl.NewDecode(bs)
	logs.PanicErr(err)

	m, _, err := d.Do(input, func(c script.Client) {
		//自定义函数
		c.SetFunc("sub0x33ReverseHEXToFloat", func(args *script.Args) (any, error) {
			msg := args.GetString(1)
			decimals := args.GetInt(2)
			bs, err := hex.DecodeString(msg)
			if err != nil {
				return nil, err
			}
			return g.Bytes(bs).Sub0x33ReverseHEXToFloat(decimals)
		})

	})
	logs.PanicErr(err)
	_ = m
	logs.Debug(m)
}
