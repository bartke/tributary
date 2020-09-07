package module

import (
	"github.com/bartke/tributary/network"
	lua "github.com/yuin/gopher-lua"
)

type Engine struct {
	network *network.Network
	exports map[string]lua.LGFunction
	vm      *VM
}

func New(n *network.Network) *Engine {
	m := &Engine{
		network: n,
	}
	m.initExports()
	return m
}

// Run will preload the tributary module and execute a script on the VM. We have to close it after
// we called Run() so we return the vm.
func (m *Engine) Run(script string) (*VM, error) {
	m.vm = NewVM()
	m.vm.LoadModule("tributary", m.Loader)
	err := m.vm.Run(script)
	return m.vm, err
}
