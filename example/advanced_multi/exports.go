package main

import (
	"strconv"

	"github.com/bartke/tributary/example/advanced/event"
	"github.com/bartke/tributary/module"
	"github.com/bartke/tributary/network"
	"github.com/bartke/tributary/pipeline/injector"
	"github.com/bartke/tributary/pipeline/gormdedupe"
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

func createDeduper(n *network.Network, f *gormdedupe.Filter) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		name := l.CheckString(1)
		d, err := f.Create(name)
		if err != nil {
			l.ArgError(1, "node not found")
			return 0
		}
		ci := interceptor.New(d.Dedupe)
		n.AddNode(name, ci)
		l.Push(module.LuaConvertValue(l, true))
		return 1
	}
}

func createTicker(n *network.Network) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		name := l.CheckString(1)
		input := l.CheckString(2)
		ms, err := strconv.Atoi(input)
		if err != nil {
			l.ArgError(1, "cannot parse ms")
			return 0
		}
		node := ticker.New(ms)
		n.AddNode(name, node)
		l.Push(module.LuaConvertValue(l, true))
		return 1
	}
}

func createCleaner(n *network.Network, f *gormdedupe.Filter) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		name := l.CheckString(1)
		table := l.CheckString(2)
		input := l.CheckString(3)
		s, err := strconv.Atoi(input)
		if err != nil {
			l.ArgError(1, "cannot parse seconds")
			return 0
		}
		node := handler.New(f.Clean(table, s))
		n.AddNode(name, node)
		l.Push(module.LuaConvertValue(l, true))
		return 1
	}
}
