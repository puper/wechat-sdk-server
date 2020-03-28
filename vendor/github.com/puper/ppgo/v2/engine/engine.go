package engine

import (
	"fmt"
	"sync"

	"github.com/puper/ppgo/v2/errors"
	"github.com/spf13/viper"
)

type Config = viper.Viper

type Builder func(engine *Engine) (interface{}, error)

type Closer interface {
	Close() error
}

func New(config *Config) *Engine {
	return &Engine{
		builders:  map[string]Builder{},
		instances: map[string]interface{}{},
		config:    config,
		graph:     newGraph(),
	}
}

type Engine struct {
	mutex     sync.RWMutex
	builders  map[string]Builder
	instances map[string]interface{}
	config    *Config
	graph     *graph
}

func (this *Engine) Register(name string, builder Builder, dependencies ...string) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if _, ok := this.builders[name]; !ok {
		this.builders[name] = builder
		this.graph.AddVertex(name)
		for _, dependency := range dependencies {
			this.graph.AddEdge(dependency, name)
		}
	}
}

func (this *Engine) Build() error {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.graph.TopologicalOrdering()
	names, err := this.graph.TopologicalOrdering()
	if err != nil {
		return err
	}
	for _, name := range names {
		if builder, ok := this.builders[name]; ok {
			var err error
			fmt.Printf("build component `%v`\n", name)
			if this.instances[name], err = builder(this); err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *Engine) Close() error {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	return this.close()
}

func (this *Engine) close() error {
	this.graph.TopologicalOrdering()
	multiErrs := errors.NewMultiErrors()
	names, err := this.graph.TopologicalOrdering()
	if err != nil {
		return multiErrs.Add(err)
	}
	for i := len(names) - 1; i >= 0; i-- {
		name := names[i]
		if instance, ok := this.instances[name]; ok {
			fmt.Printf("close component `%v`\n", name)
			delete(this.instances, name)
			if closer, ok := instance.(Closer); ok {
				multiErrs.Add(closer.Close())
			}
		}
	}
	if multiErrs.HasError() {
		return multiErrs
	}
	return nil
}

func (this *Engine) GetConfig() *Config {
	return this.config
}

func (this *Engine) Get(name string) interface{} {
	if instance, ok := this.instances[name]; ok {
		return instance
	}
	panic(fmt.Sprintf("engine: component `%v` not found", name))
}
