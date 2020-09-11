package runtime

import (
	"context"

	lua "github.com/yuin/gopher-lua"
)

type Runtime struct {
	L      *lua.LState
	cancel context.CancelFunc
}

func New() *Runtime {
	L := lua.NewState()
	ctx, cancelFn := context.WithCancel(context.Background())
	L.SetContext(ctx)
	return &Runtime{L, cancelFn}
}

func (r *Runtime) LoadModule(name string, fn func(l *lua.LState) int) {
	r.L.PreloadModule(name, fn)
}

func (r *Runtime) Compile(script string) (*lua.FunctionProto, error) {
	return compileLua(script)
}

func (r *Runtime) Execute(bc *lua.FunctionProto) error {
	return doCompiledFile(r.L, bc)
}

func (r *Runtime) Run(script string) error {
	bc, err := compileLua(script)
	if err != nil {
		return err
	}
	err = doCompiledFile(r.L, bc)
	if err != nil {
		return err
	}
	return nil
}

func (r *Runtime) Close() {
	r.cancel()
}
