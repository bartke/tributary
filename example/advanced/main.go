package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bartke/tributary"
	"github.com/bartke/tributary/event/standardevent"
	"github.com/bartke/tributary/example/advanced/event"
	"github.com/bartke/tributary/filter/gormfilter"
	"github.com/bartke/tributary/module"
	"github.com/bartke/tributary/network"
	"github.com/bartke/tributary/runtime"
	"github.com/bartke/tributary/sink/handler"
	"github.com/bartke/tributary/window/gormwindow"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	mysqlHost := "localhost"
	if len(os.Args) > 1 {
		mysqlHost = os.Args[1]
	}
	db := mysql.Open(fmt.Sprintf("root:root@tcp(%s:3306)/tb_test", mysqlHost))
	gormCfg := &gorm.Config{
		// performance
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		// isolate multiple instances on the same DB
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "example_",
		},
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
	n.AddNode("liability_printer", handler.New(out))
	n.AddNode("stake_printer", handler.New(out))
	// we can print the sources available
	fmt.Println("unconnected:")
	fmt.Println(tributary.Graphviz(n))

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

	fmt.Println("connected:")
	fmt.Println(tributary.Graphviz(n))

	n.Start()
	log.Println("running")
	<-time.After(20 * time.Second)
	log.Println("stopping")
	n.Stop()
	log.Println("stopped.")

	// blocking wait
	select {}
}

func out(e tributary.Event) {
	fmt.Println(string(e.Payload()))
}
