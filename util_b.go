package goutil

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

func ToB(b int64) (float64, string) {
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
