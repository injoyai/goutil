package script

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/injoyai/base/crypt/crc"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/net/ip"
	"github.com/injoyai/goutil/notice"
	"github.com/injoyai/goutil/oss/shell"
	"github.com/injoyai/goutil/str"
	"github.com/robertkrimen/otto"
	"math"
	"math/rand"
	"strings"
	"time"
)

var (
	holdTime  = maps.NewSafe()
	holdCount = maps.NewSafe()
)

func funcGo(args *Args) error {
	f := args.Get(1).Val().(func(i ...any) (otto.Value, error))
	go f()
	return nil
}

func funcPanic(args *Args) {
	panic(args.Get(1).Val())
}

// funcPrint 打印输出
func funcPrint(args *Args) error {
	msg := fmt.Sprint(args.Interfaces()...)
	_, err := fmt.Print(msg)
	return err
}

func funcPrintln(args *Args) error {
	msg := fmt.Sprint(args.Interfaces()...)
	_, err := fmt.Println(msg)
	return err
}

// funcPrintf 格式化打印
func funcPrintf(args *Args) error {
	list := args.Args
	msg := ""
	if len(list) > 0 {
		msg = fmt.Sprintf(list[0].String(), args.Interfaces()[1:]...)
	}
	_, err := fmt.Print(msg)
	return err
}

func funcSprintf(args *Args) any {
	list := args.Args
	if len(list) > 0 {
		return fmt.Sprintf(list[0].String(), args.Interfaces()[1:]...)
	}
	return ""
}

func funcSleep(args *Args) {
	time.Sleep(time.Millisecond * time.Duration(args.GetFloat64(1)*1000))
}

var r = rand.New(rand.NewSource(time.Now().Unix()))

func funcRand(args *Args) any {
	start := args.GetFloat64(1, 0)
	end := args.GetFloat64(2, start+1)
	decimals := args.GetInt(3, 2)
	n := r.Float64() * (end - start)
	ratio := math.Pow10(decimals)
	return float64(int64((n+start)*ratio)) / ratio
}

// funcSyncDate 从网络同步时间
func funcSyncDate(args *Args) (any, error) {
	// 阿里ntp.aliyun.com 腾讯time1.cloud.tencent.com
	host := args.GetString(1, "ntp.ntsc.ac.cn")
	result, err := shell.Exec("ntpdate " + host)
	if err != nil {
		return nil, err
	}
	return result.String(), err
}

// funcSetDate 设置时间
func funcSetDate(args *Args) (any, error) {
	dateStr := args.GetString(1, "1970-01-01 08:00:00")
	result, err := shell.Exec(fmt.Sprintf(`date --set="%s"`, dateStr))
	if err != nil {
		return nil, err
	}
	return result.String(), err
}

// funcGetJson 解析json,读取其中数据
func funcGetJson(args *Args) any {
	return conv.NewMap(args.GetString(1)).GetString(args.GetString(2))
}

// funcSpeak 播放语音
func funcSpeak(args *Args) error {
	msg := args.GetString(1)
	return notice.DefaultVoice.Speak(msg)
}

// funcBase64Encode base64编码
func funcBase64Encode(args *Args) any {
	data := args.GetString(1)
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// funcBase64Decode base64解码
func funcBase64Decode(args *Args) (any, error) {
	data := args.GetString(1)
	bs, err := base64.StdEncoding.DecodeString(data)
	return string(bs), err
}

// funcHexToBytes 字符转字节 例 "0102" >>> []byte{0x01,0x02}
func funcHexToBytes(args *Args) (any, error) {
	s := args.GetString(1)
	bs, err := hex.DecodeString(s)
	return string(bs), err
}

// funcHexToString 字节转字符 例 []byte{0x01,0x02} >>> "0102"
func funcHexToString(args *Args) any {
	return hex.EncodeToString([]byte(args.GetString(1)))
}

// funcHoldTime 连续保持时间触发
func funcHoldTime(args *Args) any {
	key := args.GetString(1)       //key(唯一标识)
	hold := args.GetBool(2)        //保持
	second := args.GetFloat64(3)   //持续时间(秒)
	reset := args.GetBool(4, true) //重置
	if hold {
		t := time.Now()
		first := holdTime.GetVar(key)
		if !first.IsNil() {
			res := first != nil && t.Sub(first.Val().(time.Time)).Seconds() > second
			if res && reset {
				holdTime.Del(key)
			}
			return res
		}
		holdTime.Set(key, t) //第一次触发
		return false
	}
	holdTime.Del(key)
	return false
}

// funcHoldCount 连续保持次数触发
func funcHoldCount(args *Args) any {
	key := args.GetString(1)       //key(唯一标识)
	rule := args.GetBool(2)        //规则
	count := args.GetInt(3)        //持续次数
	reset := args.GetBool(4, true) //重置
	if rule {
		co := holdCount.GetInt(key)
		co++
		holdCount.Set(key, co)
		res := co >= count
		if res && reset {
			holdCount.Del(key)
		}
		return res
	}
	holdCount.Del(key)
	return false
}

// funcLen 取字符长度
func funcLen(args *Args) any {
	v := args.Get(1)
	switch val := v.Val().(type) {
	case string:
		return len(val)
	case []byte:
		return len(val)
	case []any:
		return len(val)
	case map[string]any:
		return len(val)
	case map[any]any:
		return len(val)
	default:
		return len(v.String())
	}
}

// funcToInt 任意类型转int
func funcToInt(args *Args) any {
	return conv.Int64(args.GetString(1))
}

// funcToInt8 任意类型转int8
func funcToInt8(args *Args) any {
	return conv.Int8(args.GetString(1))
}

// funcToInt16 任意类型转int16
func funcToInt16(args *Args) any {
	return conv.Int16(args.GetString(1))
}

// funcToInt32 任意类型转int32
func funcToInt32(args *Args) any {
	return conv.Int32(args.GetString(1))
}

// funcToInt64 任意类型转int64
func funcToInt64(args *Args) any {
	return conv.Int64(args.GetString(1))
}

func funcToInt64Bytes(args *Args) any {
	return conv.Int64(args.GetBytes(1))
}

// funcToUint8 任意类型转uint8
func funcToUint8(args *Args) any {
	return conv.Uint8(args.GetString(1))
}

// funcToUint16 任意类型转uint8
func funcToUint16(args *Args) any {
	return conv.Uint16(args.GetString(1))
}

// funcToUint32 任意类型转uint32
func funcToUint32(args *Args) any {
	return conv.Uint32(args.GetString(1))
}

// funcToUint64 任意类型转uint64
func funcToUint64(args *Args) any {
	return conv.Uint32(args.GetString(1))
}

// funcToFloat 任意类型转浮点
func funcToFloat(args *Args) any {
	return conv.Float64(args.GetString(1))
}

// funcToFloat32 任意类型转浮点32位
func funcToFloat32(args *Args) any {
	return conv.Float32(args.GetString(1))
}

// funcToFloat64 任意类型转浮点64位
func funcToFloat64(args *Args) any {
	return conv.Float64(args.GetString(1))
}

// funcToString 任意类型转字符串
func funcToString(args *Args) any {
	return args.GetString(1)
}

// funcToBool 任意类型转bool
func funcToBool(args *Args) any {
	return conv.Bool(args.GetString(1))
}

func funcCut(args *Args) any {

	str := args.GetString(1)
	start := args.GetInt(2)
	end := args.GetInt(3)

	if end == 0 || end > len(str) {
		end = len(str)
	}
	if start < 0 {
		start = 0
	}
	if start >= len(str) {
		return ""
	}

	//字节无法映射到js
	return str[start:end]
}

// funcToBIN 数字转成2进制字符串
func funcToBIN(args *Args) any {
	v := args.Get(1)
	switch val := v.Val().(type) {
	case string:
		return conv.BINStr([]byte(val))
	}
	byte := args.GetInt(2)
	data := any(args.GetInt64(1))
	switch byte {
	case 1:
		data = conv.Uint8(data)
	case 2:
		data = conv.Uint16(data)
	case 4:
		data = conv.Uint32(data)
	case 8:
		data = conv.Uint32(data)
	default:
		data = conv.Uint16(data)
	}
	return conv.BINStr(data)
}

func funcToHex(args *Args) any {
	v := args.Get(1)
	switch val := v.Val().(type) {
	case string:
		return hex.EncodeToString([]byte(val))
	}
	data := args.GetInt64(1)
	bytes := []byte(nil)
	switch args.GetInt(2) {
	case 1:
		bytes = []byte{uint8(data)}
	case 2:
		bytes = conv.Bytes(uint16(data))
	case 4:
		bytes = conv.Bytes(uint32(data))
	case 8:
		bytes = conv.Bytes(uint64(data))
	default:
		bytes = conv.Bytes(uint8(data))
	}
	return hex.EncodeToString(bytes)
}

// funcGetByte 获取字节
func funcGetByte(args *Args) any {
	s := args.GetString(1)
	idx := args.GetInt(2)
	if len(s) > idx {
		return s[idx]
	}
	return 0
}

// funcSum 校验和
func funcSum(args *Args) any {
	sum := 0
	for _, v := range args.Args {
		sum += v.Int()
	}
	return sum
}

// funcAddInt 加减数
func funcAddInt(args *Args) any {
	s := args.GetString(1)
	add := args.GetInt(2)
	result := []byte(nil)
	for _, v := range []byte(s) {
		result = append(result, byte(int(v)+add))
	}
	return string(result)
}

// funcReverse 倒序
func funcReverse(args *Args) any {
	s := args.GetString(1)
	return str.Reverse(s)
}

// funcShell 执行脚本
func funcShell(args *Args) (any, error) {
	list := []string(nil)
	for _, v := range args.Args {
		list = append(list, v.String())
	}
	result, err := shell.Exec(list...)
	if err != nil {
		return nil, err
	}
	return result.String(), nil
}

func funcCrc16(args *Args) any {
	bs := args.Get(1).Bytes()
	table := args.GetString(2)
	param := crc.CRC16_MODBUS
	switch strings.ToUpper(table) {
	case "ARC":
		param = crc.CRC16_ARC
	case "AUG_CCITT":
		param = crc.CRC16_AUG_CCITT
	case "BUYPASS":
		param = crc.CRC16_BUYPASS
	case "CCITT_FALSE":
		param = crc.CRC16_CCITT_FALSE
	case "CDMA2000":
		param = crc.CRC16_CDMA2000
	case "DDS_110":
		param = crc.CRC16_DDS_110
	case "DECT_R":
		param = crc.CRC16_DECT_R
	case "DECT_X":
		param = crc.CRC16_DECT_X
	case "DNP":
		param = crc.CRC16_DNP
	case "EN_13757":
		param = crc.CRC16_EN_13757
	case "GENIBUS":
		param = crc.CRC16_GENIBUS
	case "MAXIM":
		param = crc.CRC16_MAXIM
	case "MCRF4XX":
		param = crc.CRC16_MCRF4XX
	case "RIELLO":
		param = crc.CRC16_RIELLO
	case "T10_DIF":
		param = crc.CRC16_T10_DIF
	case "TELEDISK":
		param = crc.CRC16_TELEDISK
	case "TMS37157":
		param = crc.CRC16_TMS37157
	case "USB":
		param = crc.CRC16_USB
	case "CRC_A":
		param = crc.CRC16_CRC_A
	case "KERMIT":
		param = crc.CRC16_KERMIT
	case "MODBUS":
		param = crc.CRC16_MODBUS
	case "X_25":
		param = crc.CRC16_X_25
	case "XMODEM":
		param = crc.CRC16_XMODEM
	}
	return crc.Encrypt16(bs, param).String()
}

func funcPing(args *Args) (any, error) {
	result, err := ip.Ping(args.GetString(1), args.Get(2).Second(1))
	return result.String(), err
}
