package handler

import (
	"github.com/bartke/tributary"
)

type Fn func(tributary.Event)

type handler struct {
	in   <-chan tributary.Event
	stop chan struct{}
	fn   Fn
}

func New(fn Fn) *handler {
	return &handler{
		fn:   fn,
		stop: make(chan struct{}),
	}
}

func (n *handler) In(ch <-chan tributary.Event) {
	n.in = ch
}

func (n *handler) Run() {
	for {
		select {
		case e := <-n.in:
			n.fn(e)
		case <-n.stop:
			return
		}
	}
}

func (n *handler) Drain() {
	close(n.stop)
}
