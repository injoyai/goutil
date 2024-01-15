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
