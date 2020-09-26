package main

import (
	"fmt"
	"log"

	"github.com/bartke/tributary"
	"github.com/bartke/tributary/module"
	"github.com/bartke/tributary/network"
	"github.com/bartke/tributary/runtime"
)

func main() {
	// create network and register nodes
	n := network.New()
	// ... add nodes for sources

	// create the tributary module that operates on the network and register the tributary module
	// exports with the runtim
	m := module.New(n)
	// ... add exports to be available in the runtime

	r := runtime.New()
	r.LoadModule(m.Loader)
	// Run will preload the tributary module and execute a script on the VM. We have to close it
	// after we called Run().
	if err := r.Run("./network.lua"); err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	// after nodes are linked up and the script is loaded without errors, run the network
	n.Start()
	fmt.Println(tributary.Graphviz(n))
	log.Println("running")

	// blocking wait
	select {}
}
