package str

// CutLeast 截取最少长度,不足最小长度,则补充string(0)
func CutLeast(str string, min int) string {
	if min < 0 {
		return ""
	}
	if len(str) < min {
		return str + string(make([]byte, min-len(str)))
	}
	return str[:min]
}
