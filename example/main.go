package main

import (
	"fmt"
	"time"
)

type Event interface {
	Payload() []byte
}

type Source interface {
	Out() chan Event
}

type Network interface {
	In(chan Event)
	Out() chan Event
}

type Sink interface {
	In(chan Event)
}

type ticker struct {
	ticker *time.Ticker
	out    chan Event
}

func NewTicker() *ticker {
	return &ticker{
		ticker: time.NewTicker(1 * time.Second),
		out:    make(chan Event),
	}
}

type timeevent struct {
	t time.Time
}

func TimeEvent(t time.Time) *timeevent {
	return &timeevent{t: t}
}

func (t timeevent) Payload() []byte {
	return []byte(t.t.String())
}

func (t *ticker) Run() {
	for {
		t.out <- TimeEvent(<-t.ticker.C)
	}
}

func (t *ticker) Out() chan Event {
	return t.out
}

type printer struct {
	in chan Event
}

func NewPrinter() *printer {
	return &printer{
		in: make(chan Event),
	}
}

func (p *printer) In(ch chan Event) {
	p.in = ch
}

func (p *printer) Run() {
	for {
		e := <-p.in
		fmt.Println(string(e.Payload()))
	}
}

func main() {
	source := NewTicker()
	go source.Run()

	sink := NewPrinter()
	sink.In(source.Out())
	go sink.Run()

	// blocking wait
	select {}
}
