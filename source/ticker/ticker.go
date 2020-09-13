package ticker

import (
	"time"

	"github.com/bartke/tributary"
	"github.com/bartke/tributary/event/timeevent"
)

type ticker struct {
	ticker *time.Ticker
	out    chan tributary.Event
}

func New(d time.Duration) *ticker {
	return &ticker{
		ticker: time.NewTicker(d),
		out:    make(chan tributary.Event),
	}
}

func (t *ticker) Run() {
	for {
		t.out <- timeevent.New(<-t.ticker.C)
	}
}

func (t *ticker) Out() <-chan tributary.Event {
	return t.out
}
