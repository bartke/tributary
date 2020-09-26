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
	m.exports = map[string]lua.LGFunction{
		"node_exists":      m.nodeExists,
		"link":             m.link,
		"fanout":           m.fanout,
		"fanin":            m.fanin,
		"create_ticker":    m.createTicker,
		"create_forwarder": m.createForwarder,
		"create_ratelimit": m.createRatelimiter,
		"create_discarder": m.createDiscarder,
		"create_tester":    m.createTester,
	}
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
