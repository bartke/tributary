package tributary

import (
	"fmt"
	"strconv"
)

const (
	graphvizHeader = `digraph G {
  rankdir=LR;
  node [shape=box, colorscheme=pastel13];
`
	graphvizFooter = `}`
)

func sourceNodes(n Network) map[string]struct{} {
	sources := map[string]struct{}{}
	for src := range n.Edges() {
		hasDest := false
		if n.NodeUnconnected(src) {
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

func drawGraphvizBootstrap(n Network) string {
	var nodes string = "\n"
	sources := sourceNodes(n)
	var i int
	for src := range sources {
		dest := "_" + strconv.Itoa(i)
		nodes += fmt.Sprintf("  %s -> %s\n", src, dest)
		nodes += fmt.Sprintf("  %s [shape=oval,fillcolor=2,style=radial];\n", src)
		nodes += fmt.Sprintf("  %s [shape=oval,fillcolor=1,style=radial,label=_];\n", dest)
		i++
	}
	for src, dests := range n.Edges() {
		if !n.NodeUnconnected(src) {
			continue
		}
		for j, dest := range dests {
			src := "_" + strconv.Itoa(i+j)
			nodes += fmt.Sprintf("  %s -> %s\n", src, dest)
			nodes += fmt.Sprintf("  %s [shape=oval,fillcolor=2,style=radial,label=_];\n", src)
			nodes += fmt.Sprintf("  %s [shape=oval,fillcolor=1,style=radial];\n", dest)
			i++
		}
	}

	return graphvizHeader + nodes + graphvizFooter
}

func drawGraphviz(n Network) string {
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

func Graphviz(n Network) string {
	if n.IsConnected() {
		return drawGraphviz(n)
	}
	return drawGraphvizBootstrap(n)
}
