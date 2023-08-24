package alg

// Point 点位
type Point struct {
	X, Y float64
}

func (this Point) Equal(p Point) bool {
	return EqualPoint(this, p)
}

func (this Point) InLine(l Line) bool {
	return PointInLine(this, l)
}

func (this Point) InRectangle(r Rectangle) bool {
	return PointInRectangle(this, r)
}

// Line 线
type Line [2]Point

func (this Line) HasPoint(p Point) bool {
	return PointInLine(p, this)
}

func (this Line) Intersect(l Line) bool {
	return LineIntersect(this, l)
}

// Rectangle 矩形
type Rectangle [2]Point

type Polygon []Point //多边形

func (this Polygon) Intersect(p Polygon) bool {
	return PolygonIntersect(this, p)
}

// EqualPoint 判断2点是否重叠
func EqualPoint(p1 Point, p2 Point) bool {
	return p1.X == p2.X && p1.Y == p2.Y
}

// PointInRectangle 判断点是否在矩形里面
func PointInRectangle(p Point, r Rectangle) bool {
	inX := (r[0].X >= p.X && p.X >= r[1].X) || (r[0].X <= p.X && p.X <= r[1].X)
	inY := (r[0].Y >= p.Y && p.Y >= r[1].Y) || (r[0].Y <= p.Y && p.Y <= r[1].Y)
	return inX && inY
}

// PointInLine 判断点是否在线上
func PointInLine(p Point, line Line) bool {
	p1, p2 := line[0], line[1]
	//判断是否矩形在范围
	if !PointInRectangle(p, Rectangle(line)) {
		return false
	}
	//判断是否是竖线
	if p1.X == p2.X {
		return true
	}
	tan1 := (p.Y - p1.Y) / (p.X - p1.X)
	tan2 := (p2.Y - p.Y) / (p2.X - p.X)
	return tan1 == tan2
}

// CrossProduct 计算两个向量的叉积
func CrossProduct(p1, p2, p3 Point) float64 {
	return (p2.X-p1.X)*(p3.Y-p1.Y) - (p2.Y-p1.Y)*(p3.X-p1.X)
}

// LineIntersect 2条线是否相交
func LineIntersect(l1, l2 Line) bool {
	cross1 := CrossProduct(l1[0], l1[1], l2[0]) * CrossProduct(l1[0], l1[1], l2[1])
	cross2 := CrossProduct(l2[0], l2[1], l1[0]) * CrossProduct(l2[0], l2[1], l1[1])
	return cross1 < 0 && cross2 < 0
}

// PolygonIntersect 判断多边形是否相交
func PolygonIntersect(p1, p2 Polygon) bool {
	/*
		射线法
		1.取多边形的任意一条变的中间点位(不能是顶点),2个点组成一条线
		2.判断该线经过的多边形的点数

		多边形相交法
		1.依次判断2个多边形的各个变是否相交

	*/

	for i := 0; i < len(p1); i++ {
		for j := 0; j < len(p2); j++ {
			s1 := Line{p1[i], p1[(i+1)%len(p1)]}
			s2 := Line{p2[j], p2[(j+1)%len(p2)]}
			if LineIntersect(s1, s2) {
				return true
			}
		}
	}
	return false
}
