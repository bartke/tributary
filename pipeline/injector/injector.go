package injector

import (
	"github.com/bartke/tributary"
)

type Fn func(tributary.Event) []tributary.Event

type injector struct {
	in   <-chan tributary.Event
	out  chan tributary.Event
	stop chan struct{}
	fn   Fn
}

func New(fn Fn) *injector {
	return &injector{
		out:  make(chan tributary.Event),
		stop: make(chan struct{}),
		fn:   fn,
	}
}

func (n *injector) In(ch <-chan tributary.Event) {
	n.in = ch
}

func (n *injector) Out() <-chan tributary.Event {
	return n.out
}

func (n *injector) Run() {
	for {
		select {
		case e := <-n.in:
			for _, msg := range n.fn(e) {
				if msg.Error() != nil {
					continue
				}
				n.out <- msg
			}
		case <-n.stop:
			return
		}
	}
}

func (n *injector) Drain() {
	close(n.stop)
}
