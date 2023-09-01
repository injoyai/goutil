package geo

type Line [2]Point

func (this Line) Equal(l Line) bool {
	return (this[0].Equal(l[0]) && this[1].Equal(l[1])) ||
		(this[0].Equal(l[1]) && this[1].Equal(l[0]))
}

func (this Line) Intersect(l Line) bool {
	return LineIntersect(this, l)
}
