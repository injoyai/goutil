package oss

import (
	"math"
	"testing"
)

func TestParseVolume(t *testing.T) {
	t.Log(ParseSize("30mb10kb").Uint64())
	t.Log(ParseSize("30mb10kb"))
	t.Log(ParseSize("20Gb10kb"))
	t.Log(ParseSize("20Gb10kb").Uint64())
	t.Log(ParseSize("20Gb10kb").SizeUnit())
	t.Log(ParseSize("1.999EB").SizeUnit())
	t.Log(Size(0))
}

func TestParseVolume2(t *testing.T) {
	t.Log(ParseSize("15.999EB").SizeUnit())
	t.Log(ParseSize("15.999EB").Uint64())
	t.Log(Size(math.MaxUint64).Uint64())
	t.Log(Size(math.MaxUint64))
	t.Log(Size(2024))
	t.Log(ParseSize("15.999xB").Uint64())
}
