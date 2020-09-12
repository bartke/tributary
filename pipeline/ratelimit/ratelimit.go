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
}

func New(atmost time.Duration) *limiter {
	return &limiter{
		out:    make(chan tributary.Event),
		last:   time.Time{},
		atmost: atmost,
	}
}

func (l *limiter) In(ch <-chan tributary.Event) {
	l.in = ch
}

func (l *limiter) Out() <-chan tributary.Event {
	return l.out
}

func (l *limiter) Run() {
	for e := range l.in {
		if time.Since(l.last) < l.atmost {
			continue
		}
		l.out <- e
		l.last = time.Now()
	}
}
