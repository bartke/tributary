package main

import (
	"fmt"
	"log"

	"github.com/bartke/tributary"
	"github.com/bartke/tributary/event/standardevent"
	"github.com/bartke/tributary/example/advanced/event"
	"github.com/bartke/tributary/filter/gormfilter"
	"github.com/bartke/tributary/module"
	"github.com/bartke/tributary/network"
	"github.com/bartke/tributary/runtime"
	"github.com/bartke/tributary/sink/handler"
	"github.com/bartke/tributary/window/gormwindow"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	file := sqlite.Open("file:test.db")
	db, err := gormwindow.New(file, &gorm.Config{}, standardevent.New, &event.Bet{}, &event.Selection{})
	if err != nil {
		log.Fatal(err)
	}

	deduper, err := gormfilter.New(file, &gorm.Config{}, standardevent.New)
	if err != nil {
		log.Fatal(err)
	}

	// create network and register nodes
	n := network.New()
	n.AddNode("streaming_ingest", NewStream())
	n.AddNode("printer", handler.New(out))
	n.AddNode("printer2", handler.New(out))

	m := module.New(n)
	m.AddWindowExports(db, &event.Bet{})
	m.AddFilterExport(deduper)

	r := runtime.New()
	r.LoadModule(m.Loader)
	err = r.Run("./network.lua")
	if err != nil {
		log.Fatal(err)
	}
	// add another script, see if it compiles first
	bc, err := r.Compile("./network2.lua")
	if err != nil {
		log.Fatal(err)
	}
	if err := r.Execute(bc); err != nil {
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
