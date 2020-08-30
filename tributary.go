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

type Sink interface {
	In(<-chan Event)
	Run()
}

type Pipeline interface {
	Source
	Sink
}

func Connect(nodeA Source, nodeB Sink) {
	nodeB.In(nodeA.Out())
}
