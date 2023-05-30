package idcard

import (
	"testing"
)

func TestNewIDCard(t *testing.T) {
	t.Log(New("330304198104049795"))
	t.Log(New("330304"))
	t.Log(New("330304202008129795").Age())
}
