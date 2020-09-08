package handler

import (
	"github.com/bartke/tributary"
)

type Fn func(tributary.Event)

type handler struct {
	in <-chan tributary.Event
	fn Fn
}

func New(fn Fn) *handler {
	return &handler{
		fn: fn,
	}
}

func (h *handler) In(ch <-chan tributary.Event) {
	h.in = ch
}

func (h *handler) Run() {
	for e := range h.in {
		h.fn(e)
	}
}
