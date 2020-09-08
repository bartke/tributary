package main

import (
	"github.com/bartke/tributary"
	"github.com/bartke/tributary/example/advanced/event"
	"github.com/bartke/tributary/module"
	"github.com/bartke/tributary/network"
	"github.com/bartke/tributary/pipeline/gormdedupe"
	"github.com/bartke/tributary/pipeline/injector"
	"github.com/bartke/tributary/pipeline/interceptor"
	"github.com/bartke/tributary/sink/handler"
	"github.com/bartke/tributary/source/ticker"
	"github.com/bartke/tributary/window"
	lua "github.com/yuin/gopher-lua"
)

// create table if not exist, cache if exists
// insert tick into table
func createWindow(n *network.Network, w window.Windower) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		name := l.CheckString(1)
		ci := interceptor.New(w.Create(&event.Bet{}))
		n.AddNode(name, ci)
		l.Push(module.LuaConvertValue(l, true))
		return 1
	}
}

// run query on table
// emit if not empty result set
func queryWindow(n *network.Network, w window.Windower) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		name := l.CheckString(1)
		query := l.CheckString(2)
		ci := injector.New(w.Query(query))
		n.AddNode(name, ci)
		l.Push(module.LuaConvertValue(l, true))
		return 1
	}
}

func createFilter(n *network.Network, f *gormdedupe.Filter) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		name := l.CheckString(1)
		seconds := l.CheckInt(2)
		deduper, err := f.Create(name)
		if err != nil {
			l.ArgError(1, "node not found")
			return 0
		}
		// add main deduper function
		ci := interceptor.New(deduper)
		n.AddNode(name, ci)
		// create cleanup routine for deduper
		src := ticker.New(seconds * 1000)
		n.AddNode(name+"_ticker", src)
		sink := handler.New(f.Clean(name, seconds))
		n.AddNode(name+"_cleaner", sink)
		tributary.Link(src, sink)
		l.Push(module.LuaConvertValue(l, true))
		return 1
	}
}
