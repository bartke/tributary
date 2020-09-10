package injector

import (
	"github.com/bartke/tributary"
)

type Fn func(tributary.Event) []tributary.Event

type injector struct {
	in  <-chan tributary.Event
	out chan tributary.Event
	fn  Fn
}

func New(fn Fn) *injector {
	return &injector{
		out: make(chan tributary.Event),
		fn:  fn,
	}
}

func (i *injector) In(ch <-chan tributary.Event) {
	i.in = ch
}

func (i *injector) Out() <-chan tributary.Event {
	return i.out
}

func (i *injector) Run() {
	for e := range i.in {
		for _, msg := range i.fn(e) {
			if msg.Error() != nil {
				continue
			}
			i.out <- msg
		}
	}
}
