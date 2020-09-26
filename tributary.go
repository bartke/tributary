package tributary

import (
	"context"
	"fmt"
	"strconv"
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

type Network interface {
	AddNode(name string, node Node) error
	GetNode(name string) (Node, bool)
	NodeExists(name string) bool
	Run()
	Edges() map[string][]string

	Link(nodeA string, nodeB string) error
	Fanin(nodeB string, nodes ...string) error
	Fanout(nodeA string, nodes ...string) error
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

const (
	graphvizHeader = `
digraph G {
  rankdir=LR;
  node [shape=box, colorscheme=pastel13];
`
	graphvizFooter = `}`
)

func sourceNodes(n Network) map[string]struct{} {
	sources := map[string]struct{}{}
	for src, _ := range n.Edges() {
		hasDest := false
		if src == "_" {
			continue
		}
		for _, dests := range n.Edges() {
			for _, dest := range dests {
				if dest == src {
					hasDest = true
				}
			}
		}
		if !hasDest {
			sources[src] = struct{}{}
		}
	}
	return sources
}

func GraphvizBootstrap(n Network) string {
	var nodes string = "\n"
	sources := sourceNodes(n)
	var i int
	for src, _ := range sources {
		dest := "_" + strconv.Itoa(i)
		nodes += fmt.Sprintf("  %s -> %s\n", src, dest)
		nodes += fmt.Sprintf("  %s [shape=oval,fillcolor=2,style=radial];\n", src)
		nodes += fmt.Sprintf("  %s [shape=oval,fillcolor=1,style=radial,label=_];\n", dest)
		i++
	}
	for src, n := range n.Edges() {
		if src != "_" {
			continue
		}
		for j, dest := range n {
			src := "_" + strconv.Itoa(i+j)
			nodes += fmt.Sprintf("  %s -> %s\n", src, dest)
			nodes += fmt.Sprintf("  %s [shape=oval,fillcolor=2,style=radial,label=_];\n", src)
			nodes += fmt.Sprintf("  %s [shape=oval,fillcolor=1,style=radial];\n", dest)
			i++
		}
	}

	return graphvizHeader + nodes + graphvizFooter
}

func Graphviz(n Network) string {
	sources := sourceNodes(n)
	var nodes string = "\n"
	for src, dests := range n.Edges() {
		for _, dest := range dests {
			nodes += fmt.Sprintf("  %s -> %s\n", src, dest)
			if _, is := sources[src]; is {
				nodes += fmt.Sprintf("  %s [shape=oval,fillcolor=2,style=radial];\n", src)
			}
			if _, is := n.Edges()[dest]; !is {
				nodes += fmt.Sprintf("  %s [shape=oval,fillcolor=1,style=radial];\n", dest)
			}
		}
	}
	return graphvizHeader + nodes + graphvizFooter
}
