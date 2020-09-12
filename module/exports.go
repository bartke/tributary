package module

import (
	"time"

	"github.com/bartke/tributary/pipeline/forwarder"
	"github.com/bartke/tributary/pipeline/ratelimit"
	"github.com/bartke/tributary/sink/handler"
	"github.com/bartke/tributary/sink/handler/discard"
	"github.com/bartke/tributary/sink/handler/tester"
	"github.com/bartke/tributary/source/ticker"
	lua "github.com/yuin/gopher-lua"
)

// register tributary functions

func (m *Engine) initExports() {
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
}

func (m *Engine) link(l *lua.LState) int {
	a := l.CheckString(1)
	b := l.CheckString(2)
	m.network.Link(a, b)
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Engine) fanout(l *lua.LState) int {
	a := l.CheckString(1)
	var dests []string
	for i := 2; i <= l.GetTop(); i++ {
		dests = append(dests, l.CheckString(i))
	}

	m.network.Fanout(a, dests...)
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Engine) fanin(l *lua.LState) int {
	b := l.CheckString(1)
	var srcs []string
	for i := 2; i <= l.GetTop(); i++ {
		srcs = append(srcs, l.CheckString(i))
	}

	m.network.Fanin(b, srcs...)
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Engine) createForwarder(l *lua.LState) int {
	name := l.CheckString(1)
	fwd := forwarder.New()
	m.network.AddNode(name, fwd)
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Engine) createRatelimiter(l *lua.LState) int {
	name := l.CheckString(1)
	atmost := l.CheckString(2)
	d, err := time.ParseDuration(atmost)
	if err != nil {
		l.ArgError(2, err.Error())
		return 0
	}
	limiter := ratelimit.New(d)
	m.network.AddNode(name, limiter)
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Engine) createTicker(l *lua.LState) int {
	name := l.CheckString(1)
	ms := l.CheckInt(2)
	limiter := ticker.New(ms)
	m.network.AddNode(name, limiter)
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Engine) createDiscarder(l *lua.LState) int {
	name := l.CheckString(1)
	sink := handler.New(discard.Handler)
	m.network.AddNode(name, sink)
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Engine) createTester(l *lua.LState) int {
	name := l.CheckString(1)
	indicator := l.CheckString(2)
	sink := handler.New(tester.New(indicator).Handler)
	m.network.AddNode(name, sink)
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Engine) nodeExists(l *lua.LState) int {
	name := l.CheckString(1)
	ok := m.network.NodeExists(name)
	if !ok {
		l.Push(LuaConvertValue(l, false))
		return 1
	}
	l.Push(LuaConvertValue(l, true))
	return 1
}
