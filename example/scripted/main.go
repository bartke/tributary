package main

import (
	"context"
	"log"

	"github.com/bartke/tributary/example/common"
	"github.com/bartke/tributary/module"
	"github.com/bartke/tributary/network"
	lua "github.com/yuin/gopher-lua"
)

func main() {
	// create network and register nodes
	n := network.New()
	n.AddNode("ticker_1s", common.NewTicker())
	n.AddNode("filter_even", common.NewFilter())
	n.AddNode("printer", common.NewPrinter())

	// create lua vm
	vm := lua.NewState()
	goCtx, ctxCancelFn := context.WithCancel(context.Background())
	vm.SetContext(goCtx)
	defer ctxCancelFn()

	// create the tributary module and register the tributary module exports
	m := module.New(n)
	vm.PreloadModule("tributary", m.Loader)
	network, err := compileLua("./example/scripted/network.lua")
	if err != nil {
		log.Fatal(err)
	}
	err = doCompiledFile(vm, network)
	if err != nil {
		log.Fatal(err)
	}

	// blocking wait
	select {}
}
