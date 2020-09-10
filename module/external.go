package module

import (
	"github.com/bartke/tributary/pipeline/injector"
	"github.com/bartke/tributary/pipeline/interceptor"
	"github.com/bartke/tributary/window"
	lua "github.com/yuin/gopher-lua"
)

func (m *Engine) ExportInterceptor(name string, ifn interceptor.Fn) {
	fn := func(l *lua.LState) int {
		name := l.CheckString(1)
		ci := interceptor.New(ifn)
		m.network.AddNode(name, ci)
		l.Push(LuaConvertValue(l, true))
		return 1
	}
	m.exports[name] = fn
}

func (m *Engine) ExportInjector(name string, ifn injector.Fn) {
	fn := func(l *lua.LState) int {
		name := l.CheckString(1)
		ci := injector.New(ifn)
		m.network.AddNode(name, ci)
		l.Push(LuaConvertValue(l, true))
		return 1
	}
	m.exports[name] = fn
}

func (m *Engine) AddWindowExports(w window.Windower, v interface{}) {
	// create table if not exist, cache if exists
	// insert tick into table
	createWindow := func(l *lua.LState) int {
		name := l.CheckString(1)
		ci := interceptor.New(w.Create(v))
		m.network.AddNode(name, ci)
		l.Push(LuaConvertValue(l, true))
		return 1
	}

	// run query on table
	// emit if not empty result set
	queryWindow := func(l *lua.LState) int {
		name := l.CheckString(1)
		query := l.CheckString(2)
		ci := injector.New(w.Query(query))
		m.network.AddNode(name, ci)
		l.Push(LuaConvertValue(l, true))
		return 1
	}

	m.Export("create_window", createWindow)
	m.Export("query_window", queryWindow)
}
