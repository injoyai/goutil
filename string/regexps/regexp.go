package regexps

import (
	"regexp"
)

const (
	Phone    = "^1[0-9]{10}$"
	Integer  = "^[0-9]+$"
	Email    = "^([a-z0-9_\\.-]+)@([\\da-z\\.-]+)\\.([a-z\\.]{2,6})$"
	Name     = `^\p{Han}{2,10}$`
	IdCard   = `^[1-9]\d{7}((0\d)|(1[0-2]))(([0|1|2]\d)|3[0-1])\d{3}$|^[1-9]\d{5}[1-9]\d{3}((0\d)|(1[0-2]))(([0|1|2]\d)|3[0-1])\d{3}([0-9]|X)$`
	Password = "^[0-9a-zA-Z@.]{6,30}$" // `(?=.*([a-zA-Z].*))(?=.*[0-9].*)[a-zA-Z0-9-*/+.~!@#$%^&*()]{6,20}$`
	LAN      = `^(\[::1\]:\d+)|(127\.0\.0\.1)|(localhost)|(10\.\d{1,3}\.\d{1,3}\.\d{1,3})|(172\.((1[6-9])|(2\d)|(3[01]))\.\d{1,3}\.\d{1,3})|(192\.168\.\d{1,3}\.\d{1,3})`
	IP       = "^(1\\d{2}|2[0-4]\\d|25[0-5]|[1-9]\\d|[1-9])\\.(1\\d{2}|2[0-4]\\d|25[0-5]|[1-9]\\d|\\d)\\.(1\\d{2}|2[0-4]\\d|25[0-5]|[1-9]\\d|\\d)\\.(1\\d{2}|2[0-4]\\d|25[0-5]|[1-9]\\d|\\d)$"
)

/************************* 常用类型 ************************/

// FindAll 提取匹配的内容
// @s ,原数据内容
// @reg,匹配规则
func FindAll(reg, s string) []string {
	sRegexp := regexp.MustCompile(reg)
	return sRegexp.FindAllString(s, -1)
}

// Is 自定义匹配
// @reg,匹配规则
// @str,匹配内容
func Is(reg string, str ...string) bool {
	var b bool
	for _, s := range str {
		b, _ = regexp.MatchString(reg, s)
		if false == b {
			return b
		}
	}
	return b
}

// FindAndReplace 匹配并替换
func FindAndReplace(reg, src, repl string) {
	sRegexp := regexp.MustCompile(reg)
	sRegexp.ReplaceAllString(src, repl)
}

// IsLAN 是否是本地局域网
func IsLAN(str ...string) bool { return Is(LAN, str...) }

// IsIP 是否是ip地址
func IsIP(str ...string) bool { return Is(IP, str...) }

// IsPhone 手提电话（不带前缀）最高11位
func IsPhone(str ...string) bool { return Is(Phone, str...) }

// IsInteger 纯整数
func IsInteger(str ...string) bool { return Is(Integer, str...) }

// IsEmail 邮箱 最高30位
func IsEmail(str ...string) bool { return Is(Email, str...) }

// IsIdCard 身份证
func IsIdCard(str ...string) bool { return Is(IdCard, str...) }

// IsName 判断是不是2-10个中文,待改进(姓氏是否正确什么的)
func IsName(str ...string) bool { return Is(Name, str...) }

/************************* 组合类型 ************************/

// IsID 数字+字母  不限制大小写 6~30位
func IsID(str ...string) bool { return Is("^[0-9a-zA-Z]{6,30}$", str...) }

// IsPwd 数字+字母+符号 6~30位
func IsPwd(str ...string) bool { return Is(Password, str...) }

/************************* 数字类型 ************************/

// IsDecimals 纯小数
func IsDecimals(str ...string) bool { return Is("^\\d+\\.[0-9]+$", str...) }

// IsTelephone 家用电话（不带前缀） 最高8位
func IsTelephone(str ...string) bool { return Is("^[0-9]{8}$", str...) }

/************************* 英文类型 *************************/

// IsEnglishLow 仅小写
func IsEnglishLow(str ...string) bool { return Is("^[a-z]+$", str...) }

// IsEnglishCap 仅大写
func IsEnglishCap(str ...string) bool { return Is("^[A-Z]+$", str...) }

// IsEnglish 大小写混合
func IsEnglish(str ...string) bool { return Is("^[A-Za-z]+$", str...) }
