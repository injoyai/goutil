package excel

import (
	"github.com/injoyai/base/str"
	"strconv"
)

// ToString 根据位置获取相应的字母,待优化,结果是从1开始
//@i,横坐标
//@v,纵坐标
func ToString(i, v int) string {
	//判断参数是否正确
	if i <= 0 || v <= 0 {
		return ""
	}

	var s string
	a := i

	for {
		m := a % 26
		if m == 0 {
			s += "Z"
		} else {
			s += string(byte(m%26 + 64))
		}
		a = (a - 1) / 26

		if a <= 26 {
			if a != 0 {
				s += string(byte(a + 64))
			}
			break
		}
	}

	//倒序
	s = str.Reverse(s)

	return s + strconv.Itoa(v)
}

// ToInt 根据excel表格位置获取坐标
//@s,表格的定位,比如 B3对应坐标2,3
func ToInt(s string) (int, int) {
	var (
		s1   string //横坐标de 参数
		s2   string //纵坐标de 参数
		num1 int    //横坐标
		num2 int    //纵坐标
		err  error
	)

	for i, v := range s {
		//判断参数是否正确
		if !(v >= 48 && v <= 57) && !(v >= 65 && v <= 90) {
			return 0, 0
		}

		//遍历到数字
		if v >= 48 && v <= 57 {
			s1 = s[:i]
			s2 = s[i:]

			num2, err = strconv.Atoi(s2)
			if err != nil {
				return 0, 0
			}
			break
		}
	}

	//处理横坐标
	n := len(s1) //- 1
	for i := 0; i < n; i++ {

		a := int(s1[i] - 64)
		for j := i; j < n-1; j++ {
			a *= 26
		}

		num1 += a
	}

	return num1, num2
}
