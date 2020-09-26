package interceptor

import (
	"github.com/bartke/tributary"
)

type Fn func(tributary.Event) tributary.Event

type interceptor struct {
	in   <-chan tributary.Event
	out  chan tributary.Event
	stop chan struct{}
	fn   Fn
}

func New(fn Fn) *interceptor {
	return &interceptor{
		out:  make(chan tributary.Event),
		stop: make(chan struct{}),
		fn:   fn,
	}
}

func (n *interceptor) In(ch <-chan tributary.Event) {
	n.in = ch
}

func (n *interceptor) Out() <-chan tributary.Event {
	return n.out
}

func (n *interceptor) Run() {
	for {
		select {
		case e := <-n.in:
			msg := n.fn(e)
			if msg.Error() != nil {
				continue
			}
			n.out <- msg
		case <-n.stop:
			return
		}
	}
}

func (n *interceptor) Drain() {
	close(n.stop)
}
