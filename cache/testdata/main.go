package main

import (
	"github.com/injoyai/goutil/cache"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/logs"
)

func main() {
	f := cache.NewFile("test")
	logs.Debug(f.GetString("a"))
	f.Set("a", 1)
	f.Set("b", 6)
	f.Save()

	{
		c := cache.NewCycle(10)
		for i := range g.Count(27) {
			c.Add(i)
		}
		logs.Debug(c.List())  //[17 18 19 20 21 22 23 24 25 26]
		logs.Debug(c.List(5)) // [22 23 24 25 26]
	}

	{
		c, err := cache.LoadingCycle("test3")
		if err != nil {
			logs.Err(err)
		} else {
			logs.Debug(c.List())
		}
	}

	{
		c1 := cache.NewCycle(10)
		for i := range g.Count(27) {
			c1.Add(i)
		}
		logs.PrintErr(c1.Save("test2"))
		c, err := cache.LoadingCycle("test2")
		if err != nil {
			logs.Err(err)
		} else {
			logs.Debug(c.List())  // [17 18 19 20 21 22 23 24 25 26]
			logs.Debug(c.List(6)) //[21 22 23 24 25 26]
		}
	}
}
