package network

import (
	"errors"

	"github.com/bartke/tributary"
)

var (
	ErrNodeNotFound     = errors.New("node doesn't exist")
	ErrUnknownNodeType  = errors.New("node type not recognized")
	ErrReservedNodeName = errors.New("node name is reserved")
)

const unconnected = "_"

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

func (n *Network) AddNode(name string, node tributary.Node) error {
	if name == unconnected {
		return ErrReservedNodeName
	}
	source, isSource := node.(tributary.Source)
	pipeline, isPipeline := node.(tributary.Pipeline)
	sink, isSink := node.(tributary.Sink)
	if isPipeline {
		n.edges[name] = []string{}
		n.pipelines[name] = pipeline
		return nil
	} else if isSource {
		n.edges[name] = []string{}
		n.sources[name] = source
		return nil
	} else if isSink {
		n.edges["_"] = append(n.edges["_"], name)
		n.sinks[name] = sink
		return nil
	}
	return ErrUnknownNodeType
}

func (n *Network) addEdge(a, b string) {
	if _, ok := n.edges[a]; !ok {
		n.edges[a] = []string{}
	}
	for i, e := range n.edges["_"] {
		if e == b {
			n.edges["_"] = append(n.edges["_"][:i], n.edges["_"][i+1:]...)
			break
		}
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

func (n *Network) Start() {
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

func (n *Network) Stop() {
	for _, node := range n.sources {
		node.Drain()
	}
	for _, node := range n.pipelines {
		node.Drain()
	}
	for _, node := range n.sinks {
		node.Drain()
	}
}

func (n *Network) NodeExists(a string) bool {
	_, ok := n.GetNode(a)
	return ok
}

func (n *Network) NodeUnconnected(a string) bool {
	if a == unconnected {
		return true
	}
	for _, node := range n.edges[unconnected] {
		if a == node {
			return true
		}
	}
	return false
}

func (n *Network) IsConnected() bool {
	// TODO: add more connectivity checks
	return len(n.edges[unconnected]) == 0
}

func (n *Network) Edges() map[string][]string {
	return n.edges
}
