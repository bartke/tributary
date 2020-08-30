package module

import (
	"errors"

	"github.com/bartke/tributary"
)

var (
	ErrNodeNotFound = errors.New("node doesn't exist")
)

type module struct {
	sources   map[string]tributary.Source
	pipelines map[string]tributary.Pipeline
	sinks     map[string]tributary.Sink
}

func New() *module {
	return &module{
		sources:   make(map[string]tributary.Source),
		pipelines: make(map[string]tributary.Pipeline),
		sinks:     make(map[string]tributary.Sink),
	}
}

func (m *module) RegisterSource(name string, node tributary.Source) {
	m.sources[name] = node
}

func (m *module) RegisterPipeline(name string, node tributary.Pipeline) {
	m.pipelines[name] = node
}

func (m *module) RegisterSink(name string, node tributary.Sink) {
	m.sinks[name] = node
}

func (m *module) GetSource(a string) (tributary.Source, error) {
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

func (m *module) GetSink(a string) (tributary.Sink, error) {
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
