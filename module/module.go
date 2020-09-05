package module

import (
	lua "github.com/yuin/gopher-lua"
)

type Module struct {
	network *Network

	exports map[string]lua.LGFunction
}

func NewModule(n *Network) *Module {
	m := &Module{
		network: n,
	}
	m.initExports()
	return m
}
