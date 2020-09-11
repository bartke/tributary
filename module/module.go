package module

import (
	"github.com/bartke/tributary/network"
	lua "github.com/yuin/gopher-lua"
)

type Engine struct {
	network *network.Network
	exports map[string]lua.LGFunction
}

func New(n *network.Network) *Engine {
	m := &Engine{
		network: n,
	}
	m.initExports()
	return m
}

func (m *Engine) Export(name string, fn lua.LGFunction) {
	m.exports[name] = fn
}

func (m *Engine) Loader(l *lua.LState) int {
	mod := l.SetFuncs(l.NewTable(), m.exports)

	l.Push(mod)
	return 1
}
