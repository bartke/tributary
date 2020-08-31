package main

import (
	"github.com/bartke/tributary"
	"github.com/bartke/tributary/example/common"
)

func main() {
	var source tributary.Source
	var pipeline tributary.Pipeline
	var sink tributary.Sink

	source = common.NewTicker()
	go source.Run()

	pipeline = common.NewFilter()
	tributary.Link(source, pipeline)
	go pipeline.Run()

	sink = common.NewPrinter()
	tributary.Link(pipeline, sink)
	go sink.Run()

	// blocking wait
	select {}
}
