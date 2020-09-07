package interceptor

import (
	"github.com/bartke/tributary"
)

type Fn func(tributary.Event) (tributary.Event, error)

type interceptor struct {
	in  <-chan tributary.Event
	out chan tributary.Event
	fn  Fn
}

func New(fn Fn) *interceptor {
	return &interceptor{
		out: make(chan tributary.Event),
		fn:  fn,
	}
}

func (i *interceptor) In(ch <-chan tributary.Event) {
	i.in = ch
}

func (i *interceptor) Out() <-chan tributary.Event {
	return i.out
}

func (i *interceptor) Run() {
	for e := range i.in {
		new, err := i.fn(e)
		if err != nil {
			continue
		}
		i.out <- new
	}
}
