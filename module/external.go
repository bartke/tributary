package module

import (
	"fmt"
	"time"

	"github.com/bartke/tributary/filter"
	"github.com/bartke/tributary/pipeline/forwarder"
	"github.com/bartke/tributary/pipeline/injector"
	"github.com/bartke/tributary/pipeline/interceptor"
	"github.com/bartke/tributary/sink/handler"
	"github.com/bartke/tributary/sink/handler/tester"
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

	// sliding_window
	// with cleanup
	slidingWindow := func(l *lua.LState) int {
		name := l.CheckString(1)
		outPort := l.CheckString(2)
		query := l.CheckString(3)
		table := l.CheckString(4)
		t := l.CheckString(5)
		timestampField := l.CheckString(6)
		d, err := time.ParseDuration(t)
		if err != nil {
			l.ArgError(3, err.Error())
			return 0
		}
		// create window consumer
		ci := interceptor.New(w.Create(v))
		createPort := name + "_create"
		m.network.AddNode(createPort, ci)
		// main query
		qi := injector.New(w.Query(query))
		queryPort := name + "_query"
		m.network.AddNode(queryPort, qi)
		// cleanup
		cleanupQuery := fmt.Sprintf(`delete from %s where  %s < %d`, table, timestampField, time.Now().Unix()-int64(d.Seconds()))
		qc := injector.New(w.Query(cleanupQuery))
		cleanupPort := name + "_window_cleanup"
		m.network.AddNode(cleanupPort, qc)
		// network ports
		fwdIn := forwarder.New()
		m.network.AddNode(name, fwdIn)
		fwdOut := forwarder.New()
		m.network.AddNode(outPort, fwdOut)
		fwdSplitter := forwarder.New()
		splitterPort := name + "_splitter"
		m.network.AddNode(splitterPort, fwdSplitter)
		limiterPort := fmt.Sprintf("%s_rate_limit_%s", name, t)
		m.addRateLimiter(limiterPort, t) // err checked above
		sink := handler.New(tester.New("-").Handler)
		sinkPort := name + "_null"
		m.network.AddNode(sinkPort, sink)
		// link up network
		m.network.Fanout(name, splitterPort, createPort)
		m.network.Link(createPort, queryPort)
		m.network.Link(queryPort, outPort)
		m.network.Link(splitterPort, limiterPort)
		m.network.Link(limiterPort, cleanupPort)
		m.network.Link(cleanupPort, sinkPort)
		l.Push(LuaConvertValue(l, true))
		return 1
	}

	// sliding_window_time, name, out, 'select ...', 'bets', '10s', 'create_time'
	m.Export("sliding_window_time", slidingWindow)
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
		tickerPort := name + "_ticker"
		m.network.AddNode(tickerPort, src)
		sink := handler.New(f.Clean(name, d))
		cleanerPort := name + "_cleaner"
		m.network.AddNode(cleanerPort, sink)
		m.network.Link(tickerPort, cleanerPort)
		l.Push(LuaConvertValue(l, true))
		return 1
	}

	m.Export("create_filter", createFilter)
}
