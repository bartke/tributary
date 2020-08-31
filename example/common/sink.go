package common

import (
	"fmt"

	"github.com/bartke/tributary"
)

type printer struct {
	in <-chan tributary.Event
}

func NewPrinter() *printer {
	return &printer{}
}

func (p *printer) In(ch <-chan tributary.Event) {
	p.in = ch
}

func (p *printer) Run() {
	for e := range p.in {
		fmt.Println(string(e.Payload()))
	}
}
