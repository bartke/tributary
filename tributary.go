package tributary

import (
	"context"
)

type Event interface {
	Context() context.Context
	Payload() []byte
	Error() error
}

type EventConstructor func(ctx context.Context, msg []byte, err error) Event

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

func Link(nodeA Source, nodeB Sink) {
	nodeB.In(nodeA.Out())
}

func Fanin(nodeB Sink, nodes ...Source) {
	out := make(chan Event)
	nodeB.In(out)

	for _, node := range nodes {
		src := node.Out()
		go func(ch <-chan Event) {
			for e := range ch {
				out <- e
			}
		}(src)
	}
}

// Fanout links the source to all outputs
func Fanout(nodeA Source, nodes ...Sink) {
	in := nodeA.Out()
	out := make([]chan Event, len(nodes))
	for i := range nodes {
		out[i] = make(chan Event)
		nodes[i].In(out[i])
	}
	go func() {
		for e := range in {
			for i := range out {
				// blocking forwarder, let's require an intact network
				out[i] <- e
				// debugging
				//go func(ch chan Event) {
				//	fmt.Println("sending", i, ch)
				//	ch <- e
				//}(out[i])
			}
		}
	}()
}
