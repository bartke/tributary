package main

import (
	"github.com/bartke/tributary"
	"github.com/bartke/tributary/module"
	"github.com/bartke/tributary/network"
	"github.com/bartke/tributary/pipeline/gormdedupe"
	"github.com/bartke/tributary/pipeline/interceptor"
	"github.com/bartke/tributary/sink/handler"
	"github.com/bartke/tributary/source/ticker"
	lua "github.com/yuin/gopher-lua"
)

func createFilter(n *network.Network, f *gormdedupe.Filter) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		name := l.CheckString(1)
		seconds := l.CheckInt(2)
		filter, err := f.Create(name)
		if err != nil {
			l.ArgError(1, "node not found")
			return 0
		}
		// add main filter function
		ci := interceptor.New(filter)
		n.AddNode(name, ci)
		// create cleanup routine for filter
		src := ticker.New(seconds * 1000)
		n.AddNode(name+"_ticker", src)
		sink := handler.New(f.Clean(name, seconds))
		n.AddNode(name+"_cleaner", sink)
		tributary.Link(src, sink)
		l.Push(module.LuaConvertValue(l, true))
		return 1
	}
}
