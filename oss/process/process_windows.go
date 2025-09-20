package process

func NewCmd(args ...string) Process {
	args = append([]string{"/c"}, args...)
	return New("cmd", args...)
}

func NewPowershell(args ...string) Process {
	args = append([]string{"/c"}, args...)
	return New("powershell", args...)
}
