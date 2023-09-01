package geo

type Polygon []Point

// Intersect 判断2个多边形是否相交,注意不包括包含关系
func (this Polygon) Intersect(p Polygon) bool {
	return PolygonIntersect(this, p)
}

// Common 获取2个多边形的相交部分
func (this Polygon) Common(p Polygon) []Polygon {
	return PolygonCommon(this, p)
}
