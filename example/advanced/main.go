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
	n.AddNode("streaming_ingest", NewStream())
	n.AddNode("printer", common.NewPrinter())

	m := module.New(n)
	m.Export("parse_json", parseJson(n))
	m.Export("create_window", parseJson(n))
	m.Export("query_window", parseJson(n))

	vm, err := m.Run("./example/advanced/network.lua")
	if err != nil {
		log.Fatal(err)
	}
	defer vm.Close()

	log.Println("running")

	// blocking wait
	select {}
}
