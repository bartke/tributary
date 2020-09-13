package module

import (
	"time"

	"github.com/bartke/tributary"
	"github.com/bartke/tributary/filter"
	"github.com/bartke/tributary/pipeline/injector"
	"github.com/bartke/tributary/pipeline/interceptor"
	"github.com/bartke/tributary/sink/handler"
	"github.com/bartke/tributary/source/ticker"
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
	m.ExportInterceptor("create_window", w.Create(v))

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

	m.Export("query_window", queryWindow)
}

func (m *Engine) AddFilterExport(f filter.Filter) {
	createFilter := func(l *lua.LState) int {
		name := l.CheckString(1)
		interval := l.CheckString(2)
		d, err := time.ParseDuration(interval)
		if err != nil {
			l.ArgError(2, err.Error())
			return 0
		}
		filter, err := f.Create(name)
		if err != nil {
			l.ArgError(1, err.Error())
			return 0
		}
		// add main filter function
		ci := interceptor.New(filter)
		m.network.AddNode(name, ci)
		// create cleanup routine for filter
		src := ticker.New(d)
		m.network.AddNode(name+"_ticker", src)
		sink := handler.New(f.Clean(name, d))
		m.network.AddNode(name+"_cleaner", sink)
		tributary.Link(src, sink)
		l.Push(LuaConvertValue(l, true))
		return 1
	}

	m.Export("create_filter", createFilter)
}
