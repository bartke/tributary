package module

import (
	"time"

	"github.com/bartke/tributary"
	"github.com/bartke/tributary/pipeline/forwarder"
	"github.com/bartke/tributary/pipeline/ratelimit"
	"github.com/bartke/tributary/sink/handler"
	"github.com/bartke/tributary/sink/handler/discard"
	"github.com/bartke/tributary/sink/handler/tester"
	lua "github.com/yuin/gopher-lua"
)

// register tributary functions

func (m *Engine) initExports() {
	m.exports = map[string]lua.LGFunction{
		"node_exists":      m.nodeExists,
		"run":              m.run,
		"link":             m.link,
		"fanout":           m.fanout,
		"fanin":            m.fanin,
		"create_forwarder": m.createForwarder,
		"create_ratelimit": m.createRatelimiter,
		"create_discarder": m.createDiscarder,
		"create_tester":    m.createTester,
	}
}

func (m *Engine) link(l *lua.LState) int {
	a, err := m.network.GetSource(l.CheckString(1))
	if err != nil {
		l.ArgError(1, "expects string "+err.Error())
		return 0
	}
	b, err := m.network.GetSink(l.CheckString(2))
	if err != nil {
		l.ArgError(2, "expects string "+err.Error())
		return 0
	}

	tributary.Link(a, b)
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Engine) fanout(l *lua.LState) int {
	src, err := m.network.GetSource(l.CheckString(1))
	if err != nil {
		l.ArgError(1, "expects string "+err.Error())
		return 0
	}
	var dests []tributary.Sink
	for i := 2; i <= l.GetTop(); i++ {
		dest, err := m.network.GetSink(l.CheckString(i))
		if err != nil {
			l.ArgError(i, "expects string "+err.Error())
			return 0
		}
		dests = append(dests, dest)
	}

	tributary.Fanout(src, dests...)
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Engine) fanin(l *lua.LState) int {
	dest, err := m.network.GetSink(l.CheckString(1))
	if err != nil {
		l.ArgError(1, "expects string "+err.Error())
		return 0
	}
	var srcs []tributary.Source
	for i := 2; i <= l.GetTop(); i++ {
		src, err := m.network.GetSource(l.CheckString(i))
		if err != nil {
			l.ArgError(i, "expects string "+err.Error())
			return 0
		}
		srcs = append(srcs, src)
	}

	tributary.Fanin(dest, srcs...)
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

func (m *Engine) run(l *lua.LState) int {
	name := l.CheckString(1)
	err := m.network.RunNode(name)
	if err != nil {
		l.ArgError(1, "node not found "+err.Error())
		return 0
	}
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
