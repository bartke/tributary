package forwarder

import (
	"github.com/bartke/tributary"
)

type forwarder struct {
	in  <-chan tributary.Event
	out chan tributary.Event
}

func New() *forwarder {
	return &forwarder{
		out: make(chan tributary.Event),
	}
}

func (f *forwarder) In(ch <-chan tributary.Event) {
	f.in = ch
}

func (f *forwarder) Out() <-chan tributary.Event {
	return f.out
}

func (f *forwarder) Run() {
	for e := range f.in {
		f.out <- e
	}
}
