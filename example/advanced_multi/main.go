package main

import (
	"fmt"
	"log"

	"github.com/bartke/tributary"
	"github.com/bartke/tributary/event/standardevent"
	"github.com/bartke/tributary/example/advanced_multi/event"
	"github.com/bartke/tributary/filter/gormfilter"
	"github.com/bartke/tributary/module"
	"github.com/bartke/tributary/network"
	"github.com/bartke/tributary/runtime"
	"github.com/bartke/tributary/sink/handler"
	"github.com/bartke/tributary/window/gormwindow"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := mysql.Open("root:root@tcp(localhost:3306)/tb_test")
	gormCfg := &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}
	window, err := gormwindow.New(db, gormCfg, standardevent.New,
		&event.Bet{},
		&event.Selection{})
	if err != nil {
		log.Fatal(err)
	}

	deduper, err := gormfilter.New(db, gormCfg, standardevent.New)
	if err != nil {
		log.Fatal(err)
	}

	// create network and register nodes
	n := network.New()
	n.AddNode("streaming_ingest", NewStream())
	n.AddNode("printer", handler.New(out))
	n.AddNode("printer2", handler.New(out))

	m := module.New(n)
	m.AddWindowExports(window, &event.Bet{})
	m.AddFilterExport(deduper)

	r := runtime.New()
	r.LoadModule(m.Loader)
	err = r.Run("./network.lua")
	if err != nil {
		log.Fatal(err)
	}
	// add another script, see if it compiles first. Note that this script depends on the first and
	// can only be loaded after. Otherwise it will error out.
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
