package module

import (
	"github.com/bartke/tributary/pipeline/injector"
	"github.com/bartke/tributary/pipeline/interceptor"
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
