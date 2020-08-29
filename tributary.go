package tributary

type Event interface {
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
