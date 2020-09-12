package network

import (
	"fmt"

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

func (n *Network) Graphviz() string {
	header := `
digraph G {
	rankdir=LR;
	node [shape=box, colorscheme=pastel13];
`
	footer := `}`

	var nodes string = "\n"
	for src, dests := range n.edges {
		for _, dest := range dests {
			nodes += fmt.Sprintf("        %s -> %s\n", src, dest)
			if _, is := n.sources[src]; is {
				nodes += fmt.Sprintf("        %s [shape=oval,fillcolor=2,style=radial];\n", src)
			}
			if _, is := n.sinks[dest]; is {
				nodes += fmt.Sprintf("        %s [shape=oval,fillcolor=1,style=radial];\n", dest)
			}
		}
	}
	return header + nodes + footer
}
