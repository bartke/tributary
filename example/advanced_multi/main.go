package main

import (
	"log"

	"github.com/bartke/tributary/example/advanced/event"
	"github.com/bartke/tributary/example/common"
	"github.com/bartke/tributary/module"
	"github.com/bartke/tributary/network"
	"github.com/bartke/tributary/pipeline/gormdedupe"
	"github.com/bartke/tributary/window/gormwindow"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	file := sqlite.Open("file:test.db")
	db, err := gormwindow.New(file, &gorm.Config{}, Msg, &event.Bet{}, &event.Selection{})
	if err != nil {
		log.Fatal(err)
	}

	deduper, err := gormdedupe.New(file, &gorm.Config{}, Msg)
	if err != nil {
		log.Fatal(err)
	}

	// create network and register nodes
	n := network.New()
	n.AddNode("streaming_ingest", NewStream())
	n.AddNode("printer", common.NewPrinter())

	m := module.New(n)
	m.AddWindowExports(db, &event.Bet{})
	m.Export("create_filter", createFilter(n, deduper))

	vm, err := m.Run("./network.lua")
	if err != nil {
		log.Fatal(err)
	}
	defer vm.Close()

	n.Run()
	log.Println("running")

	// blocking wait
	select {}
}
