package main

import (
	"log"

	"github.com/bartke/tributary/example/common"
	"github.com/bartke/tributary/module"
	"github.com/bartke/tributary/network"
	"github.com/bartke/tributary/runtime"
)

func main() {
	// create network and register nodes
	n := network.New()
	n.AddNode("ticker_1s", common.NewTicker())
	n.AddNode("filter_even", common.NewFilter())
	n.AddNode("printer", common.NewPrinter())

	// create the tributary module and register the tributary module exports
	m := module.New(n)
	r := runtime.New()
	r.LoadModule(m.Loader)
	// Run will preload the tributary module and execute a script on the VM. We have to close it
	// after we called Run().
	err := r.Run("./network.lua")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	n.Run()
	log.Println("running")

	// blocking wait
	select {}
}
