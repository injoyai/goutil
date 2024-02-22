package win

import (
	"testing"
)

func TestRegister(t *testing.T) {
	t.Log(GetRegister(REGISTER_ROOT, "m3u8dl"))
	t.Log(GetRegister(REGISTER_ROOT, "m3u8dl/shell"))
	t.Log(GetRegister(REGISTER_ROOT, "m3u8dl/shell/open/command"))
}

func TestAPPPath(t *testing.T) {
	t.Log(APPPath("chrome.exe"))
}
