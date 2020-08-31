package common

import (
	"time"

	"github.com/bartke/tributary"
)

type filter struct {
	in  <-chan tributary.Event
	out chan tributary.Event
}

func NewFilter() *filter {
	return &filter{
		out: make(chan tributary.Event),
	}
}

func (f *filter) In(ch <-chan tributary.Event) {
	f.in = ch
}

func (f *filter) Out() <-chan tributary.Event {
	return f.out
}

func (f *filter) Run() {
	for {
		e := <-f.in
		t, _ := time.Parse(time.RFC3339, string(e.Payload()))
		if t.Second()%2 == 0 {
			continue
		}
		f.out <- e
	}
}
