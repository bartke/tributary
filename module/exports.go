package module

import (
	"time"

	"github.com/bartke/tributary"
	"github.com/bartke/tributary/pipeline/forwarder"
	lua "github.com/yuin/gopher-lua"
)

// register tributary functions

func (m *Module) link(l *lua.LState) int {
	a, err := m.GetSource(l.CheckString(1))
	if err != nil {
		l.ArgError(1, "expects string")
		return 0
	}
	b, err := m.GetSink(l.CheckString(2))
	if err != nil {
		l.ArgError(2, "expects string")
		return 0
	}

	tributary.Link(a, b)
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Module) fanout(l *lua.LState) int {
	src, err := m.GetSource(l.CheckString(1))
	if err != nil {
		l.ArgError(1, "expects string")
		return 0
	}
	var dests []tributary.Sink
	for i := 2; i <= l.GetTop(); i++ {
		dest, err := m.GetSink(l.CheckString(i))
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
	dest, err := m.GetSink(l.CheckString(1))
	if err != nil {
		l.ArgError(1, "expects string")
		return 0
	}
	var srcs []tributary.Source
	for i := 2; i <= l.GetTop(); i++ {
		src, err := m.GetSource(l.CheckString(i))
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
	m.AddNode(name, fwd)
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Module) run(l *lua.LState) int {
	name := l.CheckString(1)
	err := m.Run(name)
	if err != nil {
		l.ArgError(1, "not found")
		return 0
	}
	l.Push(LuaConvertValue(l, true))
	return 1
}

func (m *Module) nodeExists(l *lua.LState) int {
	name := l.CheckString(1)
	ok := m.NodeExists(name)
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

func LuaConvertValue(l *lua.LState, val interface{}) lua.LValue {
	if val == nil {
		return lua.LNil
	}

	switch v := val.(type) {
	case bool:
		return lua.LBool(v)
	case int:
		return lua.LNumber(v)
	case int32:
		return lua.LNumber(v)
	case int64:
		return lua.LNumber(v)
	case uint32:
		return lua.LNumber(v)
	case uint64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case []byte:
		return lua.LString(v)
	case float32:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case []string:
		tbl := l.CreateTable(len(val.([]string)), 0)
		for k, v := range v {
			tbl.RawSetInt(k+1, lua.LString(v))
		}
		return tbl
	case []interface{}:
		tbl := l.CreateTable(len(val.([]interface{})), 0)
		for k, v := range v {
			tbl.RawSetInt(k+1, LuaConvertValue(l, v))
		}
		return tbl
	case time.Time:
		return lua.LNumber(v.UTC().Unix())
	default:
		return lua.LNil
	}
}
