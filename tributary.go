package tributary

import "context"

type Event interface {
	Context() context.Context
	Payload() []byte
}

type Node interface {
	Run()
}

type Source interface {
	Node
	Out() <-chan Event
}

type Sink interface {
	Node
	In(<-chan Event)
}

type Pipeline interface {
	Source
	Sink
}

func Connect(nodeA Source, nodeB Sink) {
	nodeB.In(nodeA.Out())
}
