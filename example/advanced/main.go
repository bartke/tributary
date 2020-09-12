package main

import (
	"fmt"
	"log"

	"github.com/bartke/tributary"
	"github.com/bartke/tributary/event/standardevent"
	"github.com/bartke/tributary/example/advanced/event"
	"github.com/bartke/tributary/module"
	"github.com/bartke/tributary/network"
	"github.com/bartke/tributary/runtime"
	"github.com/bartke/tributary/sink/handler"
	"github.com/bartke/tributary/window/gormwindow"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	inmemory := sqlite.Open("file::memory:?cache=shared")
	db, err := gormwindow.New(inmemory, &gorm.Config{}, standardevent.New,
		&event.Bet{},
		&event.Selection{})
	if err != nil {
		log.Fatal(err)
	}

	// create network and register nodes
	n := network.New()
	n.AddNode("streaming_ingest", NewStream())
	n.AddNode("printer", handler.New(out))

	m := module.New(n)
	m.Export("create_window", createWindow(n, db))
	m.Export("query_window", queryWindow(n, db))

	r := runtime.New()
	r.LoadModule(m.Loader)
	err = r.Run("./network.lua")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	n.Run()
	log.Println("running")

	// blocking wait
	select {}
}

func out(e tributary.Event) {
	fmt.Println(string(e.Payload()))
}
