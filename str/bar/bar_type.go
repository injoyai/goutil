package bar

type Element interface {
	String() string
}

type element func() string

func (this element) String() string { return this() }

type Format struct {
	Bar      Element
	Rate     Element
	Size     Element
	SizeUnit Element
	Speed    Element
	Used     Element
	Remain   Element
}

type Formatter func(e *Format) string

func ToB(b int64) (float64, string) {
	var mapB = map[int]string{
		0:   "B",
		10:  "KB",
		20:  "MB",
		30:  "GB",
		40:  "TB",
		50:  "PB",
		60:  "EB",
		70:  "ZB",
		80:  "YB",
		90:  "BB",
		100: "NB",
		110: "DB",
		120: "CB",
		130: "XB",
	}

	for n := 0; n <= 130; n += 10 {
		if b < 1<<(n+10) {
			if n == 0 {
				return float64(b), mapB[n]
			}
			return float64(b) / float64(int64(1)<<n), mapB[n]
		}
	}
	return float64(b), mapB[0]
}
