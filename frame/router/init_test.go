package router

import (
	"io/ioutil"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	s := New()
	s.ALL("/api/test2/1111111111111111111111111111111111", func(r *Request) {
		panic(777)
		r.WriteJson("succ")
		r.Exit()
		log.Println(666)
	})

	s.PUT("/api/test2/1", succ)
	s.PUT("/api/test2/*", func(r *Request) { r.WriteString("xxx") })

	s.Group("/api", func(g *Group) {
		g.DELETE("/test2", func(r *Request) { r.WriteString("delete") })
		g.POST("/test2", func(r *Request) { r.WriteString("post") })
		g.PUT("/test2", func(r *Request) { r.WriteString("post") })
		g.GET("/test2", func(r *Request) {
			ioutil.ReadAll(r.Request.Body)
			r.Request.Body.Close()
			r.WriteString("test")
			r.Exit()
		})
		g.Group("/api2", func(g *Group) {
			g.POST("/test2", nil)
			g.Group("/api3", func(g *Group) {
				g.ALL("/test32", nil)
				g.ALL("/test3", func(r *Request) {
					r.WriteString("test3")
					r.WriteString("test3")
					r.Exit()
				})
				g.ALL("/test34", nil)
			})
		})
	})
	s.Run()
}

func succ(r *Request) {
	r.WriteString("succ")
}
