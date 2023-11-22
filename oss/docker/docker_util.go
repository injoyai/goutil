package docker

type PageSize struct {
	PageNum  int
	PageSize int
}

// Limit 分页
func (this PageSize) Limit(total int) (int, int) {
	if total <= 0 {
		return 0, 0
	}
	if this.PageSize <= 0 {
		return 0, total
	}
	if this.PageNum < 0 {
		this.PageNum = 0
	}

	start, end := 0, total
	if e := (this.PageNum + 1) * this.PageNum; e < end {
		end = e
	}
	if s := this.PageNum * this.PageNum; s > 0 && s < end {
		start = s
	}

	return start, end
}
