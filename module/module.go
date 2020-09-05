package module

import (
	"errors"

	"github.com/bartke/tributary"
	lua "github.com/yuin/gopher-lua"
)

var (
	ErrNodeNotFound = errors.New("node doesn't exist")
)

type Module struct {
	sources   map[string]tributary.Source
	pipelines map[string]tributary.Pipeline
	sinks     map[string]tributary.Sink

	exports map[string]lua.LGFunction
}

func New() *Module {
	m := &Module{
		sources:   make(map[string]tributary.Source),
		pipelines: make(map[string]tributary.Pipeline),
		sinks:     make(map[string]tributary.Sink),
	}
	m.initExports()
	return m
}

func (m *Module) AddNode(name string, node tributary.Node) {
	source, isSource := node.(tributary.Source)
	pipeline, isPipeline := node.(tributary.Pipeline)
	sink, isSink := node.(tributary.Sink)
	if isPipeline {
		m.pipelines[name] = pipeline
	} else if isSource {
		m.sources[name] = source
	} else if isSink {
		m.sinks[name] = sink
	}
}

func (m *Module) GetSource(a string) (tributary.Source, error) {
	source, ok := m.sources[a]
	if ok {
		return source, nil
	}
	node, ok := m.pipelines[a]
	if ok {
		return node, nil
	}
	return nil, ErrNodeNotFound
}

func (m *Module) GetSink(a string) (tributary.Sink, error) {
	sink, ok := m.sinks[a]
	if ok {
		return sink, nil
	}
	node, ok := m.pipelines[a]
	if ok {
		return node, nil
	}
	return nil, ErrNodeNotFound
}

func (m *Module) GetNode(a string) (tributary.Node, bool) {
	source, ok := m.sources[a]
	if ok {
		return source, true
	}
	sink, ok := m.sinks[a]
	if ok {
		return sink, true
	}
	node, ok := m.pipelines[a]
	if ok {
		return node, true
	}
	return nil, false
}

func (m *Module) Run(a string) error {
	node, ok := m.GetNode(a)
	if ok {
		go node.Run()
		return nil
	}
	return ErrNodeNotFound
}

func (m *Module) NodeExists(a string) bool {
	_, ok := m.GetNode(a)
	return ok
}
