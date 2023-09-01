package geo

import (
	"github.com/ctessum/polyclip-go"
	"math"
)

// EqualPoint 判断2点是否重叠
func EqualPoint(p1 Point, p2 Point) bool {
	return p1.X == p2.X && p1.Y == p2.Y
}

// Distance 计算2点的距离
func Distance(p1, p2 Point) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// PointInLine 判断点是否在线上
func PointInLine(p Point, l Line) bool {
	// 计算点到线段两个端点的距离
	dist1 := Distance(p, l[0])
	dist2 := Distance(p, l[1])
	// 计算线段的总长度
	lineLength := Distance(l[0], l[1])
	// 如果点到两个端点的距离之和等于线段的总长度，则点在线上
	return math.Abs(dist1+dist2-lineLength) < 1e-9
}

func PolygonCommon(p1, p2 Polygon) []Polygon {
	ploy1 := polyclip.Polygon{
		func() (p polyclip.Contour) {
			for _, v := range p1 {
				p = append(p, polyclip.Point(v))
			}
			return
		}(),
	}
	ploy2 := polyclip.Polygon{
		func() (p polyclip.Contour) {
			for _, v := range p2 {
				p = append(p, polyclip.Point(v))
			}
			return
		}(),
	}
	common := []Polygon(nil)
	for _, v := range ploy1.Construct(polyclip.INTERSECTION, ploy2) {
		common = append(common, func() (p Polygon) {
			for _, k := range v {
				p = append(p, Point(k))
			}
			return
		}())
	}
	return common
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
	//多边形相交法
	//1.依次判断2个多边形的各个变是否相交
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
