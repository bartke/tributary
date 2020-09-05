package main

import (
	"github.com/bartke/tributary/module"
	"github.com/bartke/tributary/network"
	"github.com/bartke/tributary/pipeline/forwarder"
	lua "github.com/yuin/gopher-lua"
)

func parseJson(n *network.Network) func(l *lua.LState) int {
	return func(l *lua.LState) int {
		name := l.CheckString(1)
		fwd := forwarder.New()
		n.AddNode(name, fwd)
		l.Push(module.LuaConvertValue(l, true))
		return 1
	}
}
