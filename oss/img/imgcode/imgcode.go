package imgcode

import (
	"github.com/mojocn/base64Captcha"
)

var Default = base64Captcha.NewCaptcha(&base64Captcha.DriverString{
	Height:          80,
	Width:           240,
	NoiseCount:      50,
	ShowLineOptions: 20,
	Length:          4,
	Source:          "abcdefghjkmnpqrstuvwxyz23456789",
	Fonts:           []string{"chromohv.ttf"},
}, base64Captcha.DefaultMemStore)

// Get 获取字母数字混合验证码
func Get() (result string, base64Img string, err error) {
	return Default.Generate()
}
