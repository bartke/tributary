package multiinjector

import (
	"github.com/bartke/tributary"
)

type multiinjector struct {
	in     <-chan tributary.Event
	out    chan tributary.Event
	inject tributary.MultiInjector
}

func New(fn tributary.MultiInjector) *multiinjector {
	return &multiinjector{
		out:    make(chan tributary.Event),
		inject: fn,
	}
}

func (i *multiinjector) In(ch <-chan tributary.Event) {
	i.in = ch
}

func (i *multiinjector) Out() <-chan tributary.Event {
	return i.out
}

func (i *multiinjector) Run() {
	for e := range i.in {
		multi, err := i.inject(e)
		if err != nil {
			continue
		}
		for j := range multi {
			i.out <- multi[j]
		}
	}
}
