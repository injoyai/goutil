package regexps

import (
	"io"
	"regexp"
)

func NewFilter() *Filter {
	return &Filter{}
}

type Filter struct {
	regexp *regexp.Regexp
	enable bool
}

func (this *Filter) Enable(b ...bool) {
	this.enable = len(b) == 0 || b[0]
}

func (this *Filter) SetRegular(regular string) {
	this.regexp, _ = regexp.Compile(regular)
}

func (this *Filter) SetLike(like string) {
	this.SetRegular(".*" + like + ".*")
}

func (this *Filter) Valid(p []byte) bool {
	return !this.enable || this.regexp == nil || this.regexp.Match(p)
}

/**/

func NewFilterWriter(w io.Writer) *FilterWriter {
	return &FilterWriter{
		Writer: w,
		Filter: NewFilter(),
	}
}

type FilterWriter struct {
	io.Writer
	*Filter
}

func (this *FilterWriter) Write(p []byte) (n int, err error) {
	if this.Filter.Valid(p) {
		return this.Writer.Write(p)
	}
	return 0, nil
}
