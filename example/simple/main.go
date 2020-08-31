package main

import (
	"github.com/bartke/tributary"
	"github.com/bartke/tributary/example/common"
	"github.com/bartke/tributary/pipeline/forwarder"
)

func main() {
	var source tributary.Source
	var pipeline, fwd tributary.Pipeline
	var sink tributary.Sink

	source = common.NewTicker()
	pipeline = common.NewFilter()
	fwd = forwarder.New()
	sink = common.NewPrinter()

	//tributary.Link(source, pipeline)
	//tributary.Link(pipeline, sink)

	// from source to pipeline and forwarder
	tributary.Fanout(source, pipeline, fwd)
	// to sink from pipeline and forwarder
	tributary.Fanin(sink, pipeline, fwd)

	go source.Run()
	go pipeline.Run()
	go fwd.Run()
	go sink.Run()

	// blocking wait
	select {}
}
