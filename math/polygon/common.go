package polygon

import "math"

// TriangleOppositeSide 三角形对边
func TriangleOppositeSide(angle float64, l1, l2 float64) float64 {
	ll := l1*l1 + l2*l2 - 2*l1*l2*math.Cos(angle*math.Pi/180)
	return math.Sqrt(ll)
}
