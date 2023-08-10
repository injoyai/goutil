package router

import (
	"strings"
)

// Group 路由分组
type Group struct {
	server       *Server
	parent       []string   //父级所有
	name         string     //分组名称
	parentMiddle MiddleFunc //父级中间件
	middle       MiddleFunc //中间件
	group        []*Group   //分组
	route        *Route     //路由
}

func (g *Group) New(name string) *Group {
	return &Group{
		server:       g.server,
		parent:       append(g.parent, g.name),
		name:         name,
		parentMiddle: append(g.parentMiddle, g.middle...),
		middle:       []Handler{}, //中间件
		group:        []*Group{},  //分组
		route:        newRoute(),  //路由
	}
}

func (g *Group) ALL(pattern string, handler Handler) *Group {
	return g.Method("ALL", pattern, handler)
}

func (g *Group) GET(pattern string, handler Handler) *Group {
	return g.Method("GET", pattern, handler)
}

func (g *Group) POST(pattern string, handler Handler) *Group {
	return g.Method("POST", pattern, handler)
}

func (g *Group) PUT(pattern string, handler Handler) *Group {
	return g.Method("PUT", pattern, handler)
}

func (g *Group) DELETE(pattern string, handler Handler) *Group {
	return g.Method("DELETE", pattern, handler)
}

func (g *Group) PATCH(pattern string, handler Handler) *Group {
	return g.Method("PATCH", pattern, handler)
}

func (g *Group) HEAD(pattern string, handler Handler) *Group {
	return g.Method("HEAD", pattern, handler)
}

func (g *Group) CONNECT(pattern string, handler Handler) *Group {
	return g.Method("CONNECT", pattern, handler)
}

func (g *Group) OPTIONS(pattern string, handler Handler) *Group {
	return g.Method("OPTIONS", pattern, handler)
}

func (g *Group) TRACE(pattern string, handler Handler) *Group {
	return g.Method("TRACE", pattern, handler)
}

func (g *Group) Method(method, pattern string, handler Handler) *Group {
	return g.Group(pattern, func(g *Group) {
		g.route.Set(method, handler)
	})
}

// Group 路由分组
func (g *Group) Group(pattern string, groups ...func(group *Group)) *Group {
	g.group = append(g.group, func() (groupList []*Group) {
		for _, fn := range groups {
			for _, v := range g.group {
				if group := v.getGroup(pattern); group != nil {
					fn(group)
					return
				}
			}
			newGroup := g.New(cleanPath(pattern))
			fn(newGroup)
			groupList = append(groupList, newGroup)
		}
		return
	}()...)
	return g
}

//func (g *Group) Next(request *Request) {
//	defer func() {
//		if e := recover(); e != nil {
//			msg := fmt.Sprint(e)
//			if msg != MarkExit {
//				request.ClearBody()
//				request.WriteString(msg)
//			}
//		}
//		request.done()
//	}()
//	defer func() {
//		if e := recover(); e != nil {
//			msg := fmt.Sprint(e)
//			if msg != MarkExit {
//				request.WriteErr(errors.New(msg))
//			}
//		}
//		s.Bind.GetAndDo(request.GetStatusCode(), request)
//	}()
//	group := s.group.getGroup(request.URL.Path)
//	if group == nil {
//		request.SetStatusCode(404)
//		return
//	}
//}

func (g *Group) Do(method string, r *Request) {
	if g.route.IsEmpty() {
		r.SetStatusCode(404)
		return
	}
	handler := g.route.Get(method)
	if handler == nil {
		r.SetStatusCode(405)
		return
	}
	handler(r)
}

func (g *Group) isBottom() bool {
	return len(g.group) == 0
}

//獲取所有子集
func (g *Group) getGroupAll() []*Group {
	list := []*Group{g}
	for _, v := range g.group {
		list = append(list, v.getGroupAll()...)
	}
	return list
}

//完整请求路径
func (g *Group) fullPath() string {
	return cleanPath(strings.Join(g.parent, "") + g.name)
}

//根据完整路径获取分组信息
func (g *Group) getGroup(path string) *Group {
	name := g.name
	allPath := false
	if len(name) >= 2 && name[len(name)-2:] == "/*" {
		name = name[:len(name)-1]
		allPath = true
	}
	path = cleanPath(path)
	if len(path) > 1 {
		if n := strings.Index(path, name); n == 0 {
			path = path[len(name):]
			if len(path) == 0 || (allPath && !g.route.IsEmpty()) {
				return g
			}
			for _, v := range g.group {
				if group := v.getGroup(path); group != nil {
					return group
				}
			}
		}
	}
	return nil
}
