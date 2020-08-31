package common

import (
	"time"

	"github.com/bartke/tributary"
)

type ticker struct {
	ticker *time.Ticker
	out    chan tributary.Event
}

func NewTicker() *ticker {
	return &ticker{
		ticker: time.NewTicker(1 * time.Second),
		out:    make(chan tributary.Event),
	}
}

func (t *ticker) Run() {
	for {
		t.out <- TimeEvent(<-t.ticker.C)
	}
}

func (t *ticker) Out() <-chan tributary.Event {
	return t.out
}
