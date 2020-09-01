package main

import (
	"context"
	"log"

	"github.com/bartke/tributary/example/common"
	"github.com/bartke/tributary/module"
	lua "github.com/yuin/gopher-lua"
)

func main() {
	vm := lua.NewState()
	goCtx, ctxCancelFn := context.WithCancel(context.Background())
	vm.SetContext(goCtx)
	defer ctxCancelFn()

	m := module.New()
	// register network nodes
	m.RegisterSource("ticker_1s", common.NewTicker())
	m.RegisterPipeline("filter_even", common.NewFilter())
	m.RegisterSink("printer", common.NewPrinter())

	// register for lua 'require("tributary")
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
