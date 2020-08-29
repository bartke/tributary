package tributary

import "context"

type Event interface {
	Context() context.Context
	Payload() []byte
}

type Source interface {
	Out() <-chan Event
	Run()
}

type Network interface {
	In(<-chan Event)
	Out() <-chan Event
	Run()
}

type Sink interface {
	In(<-chan Event)
	Run()
}
