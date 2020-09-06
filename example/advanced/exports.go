package main

import (
	"github.com/bartke/tributary/module"
	"github.com/bartke/tributary/network"
	"github.com/bartke/tributary/pipeline/forwarder"
	"github.com/bartke/tributary/pipeline/injector"
	lua "github.com/yuin/gopher-lua"
)

// forward key, type, value
func parseMessage(n *network.Network) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		name := l.CheckString(1)
		fwd := forwarder.New()
		n.AddNode(name, fwd)
		l.Push(module.LuaConvertValue(l, true))
		return 1
	}
}

// create table if not exist, cache if exists
// insert tick into table
func createWindow(n *network.Network, db *Window) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		name := l.CheckString(1)
		ci := injector.New(db.createInjector)
		n.AddNode(name, ci)
		l.Push(module.LuaConvertValue(l, true))
		return 1
	}
}

// run query on table
// emit if not empty result set
func queryWindow(n *network.Network, db *Window) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		name := l.CheckString(1)
		query := l.CheckString(2)
		ci := injector.New(db.queryWindow(query))
		n.AddNode(name, ci)
		l.Push(module.LuaConvertValue(l, true))
		return 1
	}
}
