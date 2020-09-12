package tester

import (
	"fmt"

	"github.com/bartke/tributary"
)

type tester struct {
	indicator string
}

func New(indicator string) *tester {
	return &tester{indicator: indicator}
}

func (t tester) Handler(e tributary.Event) {
	if t.indicator == "" {
		fmt.Println(string(e.Payload()))
	} else {
		fmt.Print(t.indicator)
	}
}
