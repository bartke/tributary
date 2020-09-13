package network

import (
	"github.com/bartke/tributary"
)

func (n *Network) Link(a, b string) error {
	nodeA, err := n.getSource(a)
	if err != nil {
		return err
	}
	nodeB, err := n.getSink(b)
	if err != nil {
		return err
	}
	nodeB.In(nodeA.Out())
	n.addEdge(a, b)
	return nil
}

func (n *Network) Fanin(b string, an ...string) error {
	nodeB, err := n.getSink(b)
	if err != nil {
		return err
	}
	out := make(chan tributary.Event)
	nodeB.In(out)

	for _, a := range an {
		node, err := n.getSource(a)
		if err != nil {
			return err
		}
		src := node.Out()
		n.addEdge(a, b)
		go func(ch <-chan tributary.Event) {
			for e := range ch {
				out <- e
			}
		}(src)
	}
	return nil
}

// Fanout links the source to all outputs
func (n *Network) Fanout(a string, bn ...string) error {
	nodeA, err := n.getSource(a)
	if err != nil {
		return err
	}
	in := nodeA.Out()
	out := make([]chan tributary.Event, len(bn))
	for i, b := range bn {
		nodeB, err := n.getSink(b)
		if err != nil {
			return err
		}
		out[i] = make(chan tributary.Event)
		nodeB.In(out[i])
		n.addEdge(a, b)
	}
	go func() {
		for e := range in {
			for i := range out {
				out[i] <- e
			}
		}
	}()
	return nil
}
