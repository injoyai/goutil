package script

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/net/http"
	"github.com/injoyai/goutil/notice"
	"github.com/injoyai/goutil/oss/shell"
	"github.com/injoyai/goutil/str"
	"github.com/injoyai/io/dial"
	"math"
	"math/rand"
	"strings"
	"time"
)

var (
	holdTime  = maps.NewSafe()
	holdCount = maps.NewSafe()
	cacheMap  = maps.NewSafe()
	udp       = maps.NewSafe()
)

// funcPrint 打印输出
func funcPrint(args *Args) interface{} {
	msg := fmt.Sprint(args.Interfaces()...)
	fmt.Print(msg)
	return nil
}

// funcPrintf 格式化打印
func funcPrintf(args *Args) interface{} {
	list := args.Args
	msg := ""
	if len(list) > 0 {
		msg = fmt.Sprintf(list[0].String(), args.Interfaces()[1:]...)
	}
	fmt.Print(msg)
	return nil
}

func funcSprintf(args *Args) interface{} {
	list := args.Args
	if len(list) > 0 {
		return fmt.Sprintf(list[0].String(), args.Interfaces()[1:]...)
	}
	return ""
}

func funcSleep(args *Args) {
	time.Sleep(time.Duration(args.GetInt(1) * 1e6))
}

var r = rand.New(rand.NewSource(time.Now().Unix()))

func funcRand(args *Args) interface{} {
	start := args.GetFloat64(1, 0)
	end := args.GetFloat64(2, start+1)
	decimals := args.GetInt(3, 2)
	n := r.Float64() * (end - start)
	ratio := math.Pow10(decimals)
	return float64(int64((n+start)*ratio)) / ratio
}

// funcSyncDate 从网络同步时间
func funcSyncDate(args *Args) (interface{}, error) {
	// 阿里ntp.aliyun.com 腾讯time1.cloud.tencent.com
	host := args.GetString(1, "ntp.ntsc.ac.cn")
	result, err := shell.Exec("ntpdate " + host)
	return result, err
}

// funcSetDate 设置时间
func funcSetDate(args *Args) (interface{}, error) {
	dateStr := args.GetString(1, "1970-01-01 08:00:00")
	result, err := shell.Exec(fmt.Sprintf(`date --set="%s"`, dateStr))
	return result, err
}

// funcGetJson 解析json,读取其中数据
func funcGetJson(args *Args) interface{} {
	return conv.NewMap(args.GetString(1)).GetString(args.GetString(2))
}

// funcSpeak 播放语音
func funcSpeak(args *Args) error {
	msg := args.GetString(1)
	return notice.DefaultVoice.Speak(msg)
}

// funcBase64Encode base64编码
func funcBase64Encode(args *Args) interface{} {
	data := args.GetString(1)
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// funcBase64Decode base64解码
func funcBase64Decode(args *Args) (interface{}, error) {
	data := args.GetString(1)
	bs, err := base64.StdEncoding.DecodeString(data)
	return string(bs), err
}

// funcHexToBytes 字符转字节 例 "0102" >>> []byte{0x01,0x02}
func funcHexToBytes(args *Args) (interface{}, error) {
	s := args.GetString(1)
	bs, err := hex.DecodeString(s)
	return string(bs), err
}

// funcHexToString 字节转字符 例 []byte{0x01,0x02} >>> "0102"
func funcHexToString(args *Args) interface{} {
	return hex.EncodeToString([]byte(args.GetString(1)))
}

// funcHoldTime 连续保持时间触发
func funcHoldTime(args *Args) interface{} {
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
func funcHoldCount(args *Args) interface{} {
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

// funcSetCache 设置缓存
func funcSetCache(args *Args) {
	key := args.GetString(1)
	val := args.GetString(2)
	expiration := args.GetFloat64(3, 0)
	cacheMap.Set(key, val, time.Duration(float64(time.Second)*expiration))
}

// funcGetCache 获取缓存
func funcGetCache(args *Args) interface{} {
	key := args.GetString(1)
	return cacheMap.GetInterface(key)
}

// funcDelCache 删除缓存
func funcDelCache(args *Args) {
	key := args.GetString(1)
	cacheMap.Del(key)
}

// funcLen 取字符长度
func funcLen(args *Args) interface{} {
	key := args.GetString(1)
	return len(key)
}

// funcToInt 任意类型转int
func funcToInt(args *Args) interface{} {
	return conv.Int64(args.GetString(1))
}

// funcToInt8 任意类型转int8
func funcToInt8(args *Args) interface{} {
	return conv.Int8(args.GetString(1))
}

// funcToInt16 任意类型转int16
func funcToInt16(args *Args) interface{} {
	return conv.Int16(args.GetString(1))
}

// funcToInt32 任意类型转int32
func funcToInt32(args *Args) interface{} {
	return conv.Int32(args.GetString(1))
}

// funcToInt64 任意类型转int64
func funcToInt64(args *Args) interface{} {
	return conv.Int64(args.GetString(1))
}

// funcToUint8 任意类型转uint8
func funcToUint8(args *Args) interface{} {
	return conv.Uint8(args.GetString(1))
}

// funcToUint16 任意类型转uint8
func funcToUint16(args *Args) interface{} {
	return conv.Uint16(args.GetString(1))
}

// funcToUint32 任意类型转uint32
func funcToUint32(args *Args) interface{} {
	return conv.Uint32(args.GetString(1))
}

// funcToUint64 任意类型转uint64
func funcToUint64(args *Args) interface{} {
	return conv.Uint32(args.GetString(1))
}

// funcToFloat 任意类型转浮点
func funcToFloat(args *Args) interface{} {
	return conv.Float64(args.GetString(1))
}

// funcToFloat32 任意类型转浮点32位
func funcToFloat32(args *Args) interface{} {
	return conv.Float32(args.GetString(1))
}

// funcToFloat64 任意类型转浮点64位
func funcToFloat64(args *Args) interface{} {
	return conv.Float64(args.GetString(1))
}

// funcToString 任意类型转字符串
func funcToString(args *Args) interface{} {
	return args.GetString(1)
}

// funcToBool 任意类型转bool
func funcToBool(args *Args) interface{} {
	return conv.Bool(args.GetString(1))
}

// funcToBIN 数字转成2进制字符串
func funcToBIN(args *Args) interface{} {
	byte := args.GetInt(2)
	data := interface{}(args.GetInt64(1))
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

func funcToHex(args *Args) interface{} {
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
func funcGetByte(args *Args) interface{} {
	s := args.GetString(1)
	idx := args.GetInt(2)
	if len(s) > idx {
		return s[idx]
	}
	return 0
}

// funcSum 校验和
func funcSum(args *Args) interface{} {
	sum := 0
	for _, v := range args.Args {
		sum += v.Int()
	}
	return sum
}

// funcAddInt 加减数
func funcAddInt(args *Args) interface{} {
	s := args.GetString(1)
	add := args.GetInt(2)
	result := []byte(nil)
	for _, v := range []byte(s) {
		result = append(result, byte(int(v)+add))
	}
	return string(result)
}

// funcReverse 倒序
func funcReverse(args *Args) interface{} {
	s := args.GetString(1)
	return str.Reverse(s)
}

// funcShell 执行脚本
func funcShell(args *Args) (interface{}, error) {
	list := []string(nil)
	for _, v := range args.Args {
		list = append(list, v.String())
	}
	return shell.Exec(list...)
}

// funcHTTP http请求,协程执行
func funcHTTP(args *Args) error {
	method := strings.ToUpper(args.GetString(1))
	url := args.GetString(2)
	body := args.GetString(3)
	try := args.GetInt(4)
	resp := http.Url(url).SetBody(body).Retry(uint(try)).SetMethod(method).Do()
	if resp.Err() != nil {
		return resp.Err()
	}
	if conv.Int(resp.GetBodyMap()["code"]) != 200 {
		return errors.New(resp.GetBodyString())
	}
	return nil
}

func funcUDP(args *Args) error {
	addr := args.GetString(1)
	data := args.GetString(2)
	return dial.WriteUDP(addr, []byte(data))
}
