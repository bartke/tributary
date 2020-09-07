package injector

import (
	"github.com/bartke/tributary"
)

type Injector func(tributary.Event) (tributary.Event, error)

type injector struct {
	in     <-chan tributary.Event
	out    chan tributary.Event
	inject Injector
}

func New(fn Injector) *injector {
	return &injector{
		out:    make(chan tributary.Event),
		inject: fn,
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
		new, err := i.inject(e)
		if err != nil {
			continue
		}
		i.out <- new
	}
}
