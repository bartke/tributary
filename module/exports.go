package module

import (
	"time"

	"github.com/bartke/tributary"
	lua "github.com/yuin/gopher-lua"
)

// register tributary functions

func (m *module) connect(l *lua.LState) int {
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

	tributary.Connect(a, b)
	l.Push(luaConvertValue(l, true))
	return 1
}

func (m *module) Loader(l *lua.LState) int {
	exports := map[string]lua.LGFunction{
		"connect": m.connect,
	}
	mod := l.SetFuncs(l.NewTable(), exports)

	l.Push(mod)
	return 1
}

func luaConvertValue(l *lua.LState, val interface{}) lua.LValue {
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
			tbl.RawSetInt(k+1, luaConvertValue(l, v))
		}
		return tbl
	case time.Time:
		return lua.LNumber(v.UTC().Unix())
	default:
		return lua.LNil
	}
}
