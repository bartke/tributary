package main

import (
	"log"

	"github.com/bartke/tributary/example/common"
	"github.com/bartke/tributary/module"
	"github.com/bartke/tributary/network"
)

func main() {
	// create network and register nodes
	n := network.New()
	n.AddNode("ticker_1s", common.NewTicker())
	n.AddNode("filter_even", common.NewFilter())
	n.AddNode("printer", common.NewPrinter())

	// create the tributary module and register the tributary module exports
	m := module.New(n)
	vm, err := m.Run("./network.lua")
	if err != nil {
		log.Fatal(err)
	}
	defer vm.Close()

	log.Println("running")

	// blocking wait
	select {}
}
