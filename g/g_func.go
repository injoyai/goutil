package g

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/injoyai/base/chans"
	"github.com/injoyai/base/crypt/md5"
	"github.com/injoyai/base/maps/wait"
	"github.com/injoyai/conv"
	uuid "github.com/satori/go.uuid"
	"math"
	"reflect"
	"runtime"
	"runtime/debug"
	"time"
)

//========================================Chan========================================

// Range 仿python的range 参数1-3个
// 当1个参数, 例 Range(5) 输出 0,1,2,3,4
// 当2个参数, 例 Range(1,5) 输出 1,2,3,4
// 当3个参数, 例 Range(0,5,2) 输出 0,2,4
func Range(n int, m ...int) <-chan int {
	return chans.Range(n, m...)
}

// Count 遍历
// @num,数量,-1为死循环
// @interval,间隔
func Count(num int, interval ...time.Duration) <-chan int {
	return chans.Count(num, interval...)
}

// Interval 间隔触发
func Interval(interval time.Duration, nums ...int) <-chan int {
	num := conv.Default[int](-1, nums...)
	return chans.Count(num, interval)
}

//========================================Crypt========================================

// MD5 加密,返回hex的32位小写
func MD5(s string) string { return md5.Encrypt(s) }

// HmacMD5 加密,返回hex的32位小写
func HmacMD5(s, key string) string { return md5.Hmac(s, key) }

func Base64(bs []byte) string { return base64.StdEncoding.EncodeToString(bs) }

func HEX(bs []byte) string { return hex.EncodeToString(bs) }

//========================================Runtime========================================

// Recover 错误捕捉
func Recover(err *error, stack ...bool) {
	if er := recover(); er != nil {
		if err != nil {
			if len(stack) > 0 && stack[0] {
				*err = fmt.Errorf("%v\n%v", er, string(debug.Stack()))
			} else {
				*err = fmt.Errorf("%v", er)
			}
		}
	}
}

// RecoverFunc 捕捉错误并执行函数
func RecoverFunc(fn func(err error, stack string)) {
	if er := recover(); er != nil {
		if fn != nil {
			fn(fmt.Errorf("%v", er), string(debug.Stack()))
		}
	}
}

// Try 尝试运行,捕捉错误 其他语言的try catch
func Try(fn func() error, catch ...func(err error)) (err error) {
	defer RecoverFunc(func(er error, stack string) {
		err = er
		for _, v := range catch {
			v(er)
		}
	})
	return fn()
}

// Retry 重试,可选重试间隔函数,入参是0
func Retry(fn func() error, num int, interval ...func(time.Duration) time.Duration) (err error) {
	t := time.Duration(0)
	for i := 0; num < 0 || i < num; i++ {
		if err = Try(fn); err == nil {
			return
		}
		for _, v := range interval {
			t = v(t)
		}
		if t > 0 {
			<-time.After(t)
		}
	}
	return
}

// WithRetreat32 默认退避重试,最小1秒,最大32秒,
var WithRetreat32 = RetreatRange(time.Second, time.Second*32)

// RetreatRange 默认退避重试,最小1秒,最大32秒,
func RetreatRange(min, max time.Duration) func(t time.Duration) time.Duration {
	return func(t time.Duration) time.Duration {
		if t < min {
			return min
		}
		t *= 2
		if t <= max {
			return t
		}
		return max
	}
}

// PanicErr 如果是错误则panic
func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

// StopWithErr 按顺序执行函数,遇到错误结束,并返回错误
func StopWithErr(fn ...func() error) error {
	for _, v := range fn {
		if err := v(); err != nil {
			return err
		}
	}
	return nil
}

//========================================Wait========================================

// Wait 等待
func Wait(key string) (interface{}, error) { return wait.Wait(key) }

// Done 结束等待
func Done(key string, value interface{}) { wait.Done(key, value) }

//========================================OS========================================

// Input 监听用户输入
func Input(hint ...interface{}) (s string) {
	if len(hint) > 0 {
		fmt.Println(hint...)
	}
	fmt.Scanln(&s)
	return
}

// InputUntil 监听用户输入直到满足条件
func InputUntil(hint string, f func(s string) bool) (s string) {
	for {
		s = Input(hint)
		if f(s) {
			return
		}
	}
}

// InputVar 监听用户输入,返回*conv.Var
func InputVar(hint ...interface{}) *conv.Var {
	input := Input(hint...)
	if len(input) == 0 {
		return conv.Nil
	}
	return conv.New(input)
}

// FuncName 获取函数名
func FuncName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

//========================================Math========================================

// Decimals 保留小数点,默认2位
func Decimals(f float64, d ...int) float64 {
	b := math.Pow10(conv.Default[int](2, d...))
	return float64(int64(f*b)) / b
}

//========================================Third========================================

// UUID uuid
func UUID() string { return uuid.NewV4().String() }
