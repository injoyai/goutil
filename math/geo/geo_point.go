package geo

type Point struct {
	X, Y float64
}

func (this Point) Equal(p Point) bool {
	return this.X == p.X && this.Y == p.Y
}

func (this Point) Distance(p Point) float64 {
	return Distance(this, p)
}

func (this Point) InLine(l Line) bool {
	return PointInLine(this, l)
}
