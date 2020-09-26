package ratelimit

import (
	"time"

	"github.com/bartke/tributary"
)

type limiter struct {
	in     <-chan tributary.Event
	out    chan tributary.Event
	last   time.Time
	atmost time.Duration

	stop chan struct{}
}

func New(atmost time.Duration) *limiter {
	return &limiter{
		out:    make(chan tributary.Event),
		stop:   make(chan struct{}),
		last:   time.Time{},
		atmost: atmost,
	}
}

func (n *limiter) In(ch <-chan tributary.Event) {
	n.in = ch
}

func (n *limiter) Out() <-chan tributary.Event {
	return n.out
}

func (n *limiter) Run() {
	for {
		select {
		case e := <-n.in:
			if time.Since(n.last) < n.atmost {
				continue
			}
			n.out <- e
			n.last = time.Now()
		case <-n.stop:
			return
		}
	}
}

func (n *limiter) Drain() {
	close(n.stop)
}
