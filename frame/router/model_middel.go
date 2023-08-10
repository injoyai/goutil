package router

type Middle struct {
	group   *Group
	method  string
	request *Request
}

func (m *Middle) Use(handler Handler) {

}

func (m *Middle) Next() {
	m.group.Do(m.method, m.request)
}

func MiddleCORS(r *Request) {
	for key, val := range CORS {
		r.SetHeader(key, val)
	}
	r.Middle.Next()
}
