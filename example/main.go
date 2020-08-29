package main

import (
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
	t time.Time
}

func TimeEvent(t time.Time) *timeevent {
	return &timeevent{t: t}
}

func (t timeevent) Payload() []byte {
	return []byte(t.t.Format(time.RFC3339))
}

func (t *ticker) Run() {
	for {
		t.out <- TimeEvent(<-t.ticker.C)
	}
}

func (t *ticker) Out() chan tributary.Event {
	return t.out
}

type printer struct {
	in chan tributary.Event
}

func NewPrinter() *printer {
	return &printer{
		in: make(chan tributary.Event),
	}
}

func (p *printer) In(ch chan tributary.Event) {
	p.in = ch
}

func (p *printer) Run() {
	for {
		e := <-p.in
		fmt.Println(string(e.Payload()))
	}
}

type filter struct {
	in  chan tributary.Event
	out chan tributary.Event
}

func NewFilter() *filter {
	return &filter{
		out: make(chan tributary.Event),
	}
}

func (f *filter) In(ch chan tributary.Event) {
	f.in = ch
}

func (f *filter) Out() chan tributary.Event {
	return f.out
}

func (f *filter) Run() {
	for {
		e := <-f.in
		t, err := time.Parse(time.RFC3339, string(e.Payload()))
		if err != nil || t.Second()%2 == 0 {
			//fmt.Println("filtered")
			continue
		}
		f.out <- e
	}
}

func main() {
	source := NewTicker()

	pipeline := NewFilter()
	pipeline.In(source.Out())

	sink := NewPrinter()
	sink.In(pipeline.Out())

	go source.Run()
	go pipeline.Run()
	go sink.Run()

	// blocking wait
	select {}
}
