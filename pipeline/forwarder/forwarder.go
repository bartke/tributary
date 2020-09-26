package forwarder

import (
	"github.com/bartke/tributary"
)

type forwarder struct {
	in  <-chan tributary.Event
	out chan tributary.Event

	stop chan struct{}
}

func New() *forwarder {
	return &forwarder{
		out:  make(chan tributary.Event),
		stop: make(chan struct{}),
	}
}

func (n *forwarder) In(ch <-chan tributary.Event) {
	n.in = ch
}

func (n *forwarder) Out() <-chan tributary.Event {
	return n.out
}

func (n *forwarder) Run() {
	for {
		select {
		case e := <-n.in:
			n.out <- e
		case <-n.stop:
			return
		}
	}
}

func (n *forwarder) Drain() {
	close(n.stop)
}
