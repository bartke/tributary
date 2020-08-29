package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bartke/tributary"
)

type ticker struct {
	ticker *time.Ticker
	out    chan tributary.Event
}

func NewTicker() *ticker {
	return &ticker{
		ticker: time.NewTicker(1 * time.Second),
		out:    make(chan tributary.Event),
	}
}

type timeevent struct {
	t   time.Time
	ctx context.Context
}

func TimeEvent(t time.Time) *timeevent {
	return &timeevent{t: t, ctx: context.Background()}
}

func (t timeevent) Payload() []byte {
	return []byte(t.t.Format(time.RFC3339))
}

func (t timeevent) Context() context.Context {
	return t.ctx
}

func (t *ticker) Run() {
	for {
		t.out <- TimeEvent(<-t.ticker.C)
	}
}

func (t *ticker) Out() <-chan tributary.Event {
	return t.out
}

type printer struct {
	in <-chan tributary.Event
}

func NewPrinter() *printer {
	return &printer{
		in: make(chan tributary.Event),
	}
}

func (p *printer) In(ch <-chan tributary.Event) {
	p.in = ch
}

func (p *printer) Run() {
	for {
		e := <-p.in
		fmt.Println(string(e.Payload()))
	}
}

type filter struct {
	in  <-chan tributary.Event
	out chan tributary.Event
}

func NewFilter() *filter {
	return &filter{
		out: make(chan tributary.Event),
	}
}

func (f *filter) In(ch <-chan tributary.Event) {
	f.in = ch
}

func (f *filter) Out() <-chan tributary.Event {
	return f.out
}

func (f *filter) Run() {
	for {
		e := <-f.in
		t, _ := time.Parse(time.RFC3339, string(e.Payload()))
		if t.Second()%2 == 0 {
			continue
		}
		f.out <- e
	}
}

func main() {
	var source tributary.Source
	var pipeline tributary.Network
	var sink tributary.Sink

	source = NewTicker()
	go source.Run()

	pipeline = NewFilter()
	pipeline.In(source.Out())
	go pipeline.Run()

	sink = NewPrinter()
	sink.In(pipeline.Out())
	go sink.Run()

	// blocking wait
	select {}
}
