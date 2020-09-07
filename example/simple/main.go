package main

import (
	"github.com/bartke/tributary"
	"github.com/bartke/tributary/example/common"
	"github.com/bartke/tributary/pipeline/forwarder"
)

func main() {
	var source tributary.Source
	var pipeline, fwd1, fwd2 tributary.Pipeline
	var sink tributary.Sink

	source = common.NewTicker()
	pipeline = common.NewFilter()
	fwd1 = forwarder.New()
	fwd2 = forwarder.New()
	sink = common.NewPrinter()

	// from source to pipeline and forwarder via fwd1
	tributary.Link(source, fwd1)
	tributary.Fanout(fwd1, pipeline, fwd2)
	// to sink from pipeline and forwarder
	tributary.Fanin(sink, pipeline, fwd2)

	go source.Run()
	go pipeline.Run()
	go fwd1.Run()
	go fwd2.Run()
	go sink.Run()

	// blocking wait
	select {}
}
