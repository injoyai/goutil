package main

import (
	"github.com/injoyai/goutil/frame/in/v3"
	"github.com/injoyai/goutil/frame/middle"
	"github.com/injoyai/goutil/frame/mux"
)

func main() {
	s := mux.New(
		mux.WithLog(),
		mux.WithPort(8089),
		mux.WithPing(),
		mux.WithSwagger(middle.DefaultSwagger),
	)

	s.Group("/api", func(g *mux.Grouper) {
		g.ALL("/test", func(r *mux.Request) {
			in.Succ(666)
		})
	})

	s.Run()
}
