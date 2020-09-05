package module

import (
	"context"

	lua "github.com/yuin/gopher-lua"
)

type VM struct {
	L      *lua.LState
	cancel context.CancelFunc
}

func NewVM() *VM {
	L := lua.NewState()
	ctx, cancelFn := context.WithCancel(context.Background())
	L.SetContext(ctx)
	return &VM{L, cancelFn}
}

func (vm *VM) LoadModule(name string, fn func(l *lua.LState) int) {
	vm.L.PreloadModule(name, fn)
}

func (vm *VM) Run(script string) error {
	bc, err := compileLua(script)
	if err != nil {
		return err
	}
	err = doCompiledFile(vm.L, bc)
	if err != nil {
		return err
	}
	return nil
}

func (vm *VM) Close() {
	vm.cancel()
}
