package module

import (
	"time"

	lua "github.com/yuin/gopher-lua"
)

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
