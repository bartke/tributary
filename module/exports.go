package module

import (
	"github.com/bartke/tributary"
	"github.com/bartke/tributary/pipeline/forwarder"
	lua "github.com/yuin/gopher-lua"
)

// register tributary functions

func (m *Module) link(l *lua.LState) int {
	a, err := m.network.GetSource(l.CheckString(1))
	if err != nil {
		l.ArgError(1, "expects string")
		return 0
	}
	b, err := m.network.GetSink(l.CheckString(2))
	if err != nil {
		l.ArgError(2, "expects string")
		return 0
	}

	tributary.Link(a, b)
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Module) fanout(l *lua.LState) int {
	src, err := m.network.GetSource(l.CheckString(1))
	if err != nil {
		l.ArgError(1, "expects string")
		return 0
	}
	var dests []tributary.Sink
	for i := 2; i <= l.GetTop(); i++ {
		dest, err := m.network.GetSink(l.CheckString(i))
		if err != nil {
			l.ArgError(i, "expects string")
			return 0
		}
		dests = append(dests, dest)
	}

	tributary.Fanout(src, dests...)
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Module) fanin(l *lua.LState) int {
	dest, err := m.network.GetSink(l.CheckString(1))
	if err != nil {
		l.ArgError(1, "expects string")
		return 0
	}
	var srcs []tributary.Source
	for i := 2; i <= l.GetTop(); i++ {
		src, err := m.network.GetSource(l.CheckString(i))
		if err != nil {
			l.ArgError(i, "expects string")
			return 0
		}
		srcs = append(srcs, src)
	}

	tributary.Fanin(dest, srcs...)
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Module) createForwarder(l *lua.LState) int {
	name := l.CheckString(1)
	fwd := forwarder.New()
	m.network.AddNode(name, fwd)
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Module) run(l *lua.LState) int {
	name := l.CheckString(1)
	err := m.network.RunNode(name)
	if err != nil {
		l.ArgError(1, "node not found")
		return 0
	}
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Module) nodeExists(l *lua.LState) int {
	name := l.CheckString(1)
	ok := m.network.NodeExists(name)
	if !ok {
		l.Push(LuaConvertValue(l, false))
		return 1
	}
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Module) initExports() {
	m.exports = map[string]lua.LGFunction{
		"node_exists":      m.nodeExists,
		"run":              m.run,
		"link":             m.link,
		"create_forwarder": m.createForwarder,
		"fanout":           m.fanout,
		"fanin":            m.fanin,
	}
}

func (m *Module) Export(name string, fn lua.LGFunction) {
	m.exports[name] = fn
}

func (m *Module) Loader(l *lua.LState) int {
	mod := l.SetFuncs(l.NewTable(), m.exports)

	l.Push(mod)
	return 1
}
