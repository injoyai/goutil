package times

import (
	"testing"
)

func TestNow(t *testing.T) {
	now := Now()
	t.Log(now)
	t.Log(now.IntegerSec().UnixNano())
	t.Log(now.IntegerMin())
	t.Log(now.IntegerHour())
	t.Log(now.IntegerDay())
	t.Log(now.IntegerWeek())
	t.Log(now.IntegerMonth())
	t.Log(now.IntegerQuarter())
	t.Log(now.IntegerYear())
}
