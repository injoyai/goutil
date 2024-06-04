package oss

import (
	"math"
	"testing"
)

func TestParseVolume(t *testing.T) {
	t.Log(ParseVolume("30mb10kb").Uint64())
	t.Log(ParseVolume("30mb10kb"))
	t.Log(ParseVolume("20Gb10kb"))
	t.Log(ParseVolume("20Gb10kb").Uint64())
	t.Log(ParseVolume("20Gb10kb").SizeUnit())
	t.Log(ParseVolume("1.999EB").SizeUnit())
}

func TestParseVolume2(t *testing.T) {
	t.Log(ParseVolume("15.999EB").SizeUnit())
	t.Log(ParseVolume("15.999EB").Uint64())
	t.Log(Volume(math.MaxUint64).Uint64())
	t.Log(Volume(math.MaxUint64))
	t.Log(Volume(2024))
	t.Log(ParseVolume("15.999xB").Uint64())
}
