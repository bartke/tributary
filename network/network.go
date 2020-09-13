package network

import (
	"errors"

	"github.com/bartke/tributary"
)

var (
	ErrNodeNotFound = errors.New("node doesn't exist")
)

type Network struct {
	sources   map[string]tributary.Source
	pipelines map[string]tributary.Pipeline
	sinks     map[string]tributary.Sink

	edges map[string][]string
}

func New() *Network {
	n := &Network{
		sources:   make(map[string]tributary.Source),
		pipelines: make(map[string]tributary.Pipeline),
		sinks:     make(map[string]tributary.Sink),
		edges:     make(map[string][]string),
	}
	return n
}

func (n *Network) AddNode(name string, node tributary.Node) {
	source, isSource := node.(tributary.Source)
	pipeline, isPipeline := node.(tributary.Pipeline)
	sink, isSink := node.(tributary.Sink)
	if isPipeline {
		n.edges[name] = []string{}
		n.pipelines[name] = pipeline
	} else if isSource {
		n.edges[name] = []string{}
		n.sources[name] = source
	} else if isSink {
		n.sinks[name] = sink
	}
}

func (n *Network) addEdge(a, b string) {
	if _, ok := n.edges[a]; !ok {
		n.edges[a] = []string{}
	}
	n.edges[a] = append(n.edges[a], b)
}

func (n *Network) getSource(a string) (tributary.Source, error) {
	source, ok := n.sources[a]
	if ok {
		return source, nil
	}
	node, ok := n.pipelines[a]
	if ok {
		return node, nil
	}
	return nil, ErrNodeNotFound
}

func (n *Network) getSink(a string) (tributary.Sink, error) {
	sink, ok := n.sinks[a]
	if ok {
		return sink, nil
	}
	node, ok := n.pipelines[a]
	if ok {
		return node, nil
	}
	return nil, ErrNodeNotFound
}

func (n *Network) GetNode(a string) (tributary.Node, bool) {
	source, ok := n.sources[a]
	if ok {
		return source, true
	}
	sink, ok := n.sinks[a]
	if ok {
		return sink, true
	}
	node, ok := n.pipelines[a]
	if ok {
		return node, true
	}
	return nil, false
}

func (n *Network) Run() {
	for _, node := range n.sinks {
		go node.Run()
	}
	for _, node := range n.pipelines {
		go node.Run()
	}
	for _, node := range n.sources {
		go node.Run()
	}
}

func (n *Network) NodeExists(a string) bool {
	_, ok := n.GetNode(a)
	return ok
}

func (n *Network) Edges() map[string][]string {
	return n.edges
}
