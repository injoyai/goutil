package paging

func FindAndCount(i Interface, size int, offset ...int) (interface{}, int) {
	return Limit(i, size, offset...), i.Len()
}

func Limit(i Interface, size int, offset ...int) interface{} {
	if size <= 0 {
		return i.Cut(0, 0)
	}

	start, end := 0, size
	if len(offset) > 0 {
		start += offset[0]
		end += offset[0]
	}

	if start >= i.Len() {
		return i.Cut(0, 0)
	}

	if start < 0 {
		start = 0
	}
	if end >= i.Len() {
		end = i.Len()
	}

	return i.Cut(start, end)
}

type Interface interface {
	Cut(int, int) interface{}
	Len() int
}
