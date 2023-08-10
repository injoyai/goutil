package router

// New 新建服务
// 路由注册有先后顺序(注册相同的路由会执行先注册的),
// Handler是nil的时候会过滤
func New() *Server {
	return NewServer()
}
