package module

import (
	"github.com/bartke/tributary/network"
	lua "github.com/yuin/gopher-lua"
)

type Module struct {
	network *network.Network
	exports map[string]lua.LGFunction
}

func New(n *network.Network) *Module {
	m := &Module{
		network: n,
	}
	m.initExports()
	return m
}
