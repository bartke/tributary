package main

import (
	"log"

	"github.com/bartke/tributary/example/advanced/event"
	"github.com/bartke/tributary/example/common"
	"github.com/bartke/tributary/module"
	"github.com/bartke/tributary/network"
	gormwindow "github.com/bartke/tributary/window/gorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// create network and register nodes
	n := network.New()
	n.AddNode("streaming_ingest", NewStream())
	n.AddNode("printer", common.NewPrinter())

	inmemory := sqlite.Open("file::memory:?cache=shared")
	//file := sqlite.Open("file:test.db")
	db, err := gormwindow.New(inmemory, &gorm.Config{}, Msg, &event.Bet{}, &event.Selection{})
	if err != nil {
		log.Fatal(err)
	}

	m := module.New(n)
	m.Export("parse", parseMessage(n))
	m.Export("create_window", createWindow(n, db))
	m.Export("query_window", queryWindow(n, db))

	vm, err := m.Run("./example/advanced/network.lua")
	if err != nil {
		log.Fatal(err)
	}
	defer vm.Close()

	log.Println("running")

	// blocking wait
	select {}
}
