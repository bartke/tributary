package tributary

type Event interface {
	Payload() []byte
}

type Source interface {
	Out() chan Event
}

type Network interface {
	In(chan Event)
	Out() chan Event
}

type Sink interface {
	In(chan Event)
}
