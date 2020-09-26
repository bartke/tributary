package ticker

import (
	"time"

	"github.com/bartke/tributary"
	"github.com/bartke/tributary/event/timeevent"
)

type ticker struct {
	ticker *time.Ticker
	out    chan tributary.Event

	stop chan struct{}
}

func New(d time.Duration) *ticker {
	return &ticker{
		ticker: time.NewTicker(d),
		stop:   make(chan struct{}),
		out:    make(chan tributary.Event),
	}
}

func (n *ticker) Out() <-chan tributary.Event {
	return n.out
}

func (n *ticker) Run() {
	for {
		select {
		case n.out <- timeevent.New(<-n.ticker.C):
		case <-n.stop:
			return
		}
	}
}

func (n *ticker) Drain() {
	close(n.stop)
}
