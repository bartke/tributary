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
	// register sources, pipelines, sinks
	tickerNode := common.NewTicker()
	m.RegisterSource("ticker_1s", tickerNode)
	filterNode := common.NewFilter()
	m.RegisterPipeline("filter_even", filterNode)
	printerNode := common.NewPrinter()
	m.RegisterSink("printer", printerNode)

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

	go tickerNode.Run()
	go filterNode.Run()
	go printerNode.Run()

	// blocking wait
	select {}
}

var network = `
local tb = require("tributary")

tb.connect("ticker_1s", "filter_even")
tb.connect("filter_even", "printer")

`
