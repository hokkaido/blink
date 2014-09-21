package blink

import (
	"sync"
)

// A concrete Provider will be responsible for the actual retrieval of Tiles.
type Provider interface {
	GetTile(zoom int, x int, y int) ([]byte, error)
}

// Factory method for the creation of a Provider. Accepts the configuration.
type providerFunc func(interface{}) (Provider, error)

type providerConfigFunc func() interface{}

// Holds the factory methods for the creation and configuration of a provider
type providerRegistration struct {
	provider func(interface{}) (Provider, error)
	config   func() interface{}
}

type providerRegistry struct {
	registrations map[string]*providerRegistration
	mu            sync.RWMutex
}

var defaultProviderRegistry = &providerRegistry{registrations: make(map[string]*providerRegistration)}

// Registers a ProviderFunc
func (r *providerRegistry) register(name string, pr *providerRegistration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.registrations == nil {
		r.registrations = make(map[string]*providerRegistration)
	}
	if _, exists := r.registrations[name]; exists {
		panic("A provider with the name " + name + " is already registered")
	}
	r.registrations[name] = pr
}

func RegisterProvider(name string, pf func(interface{}) (Provider, error), pcf func() interface{}) {
	var registration = &providerRegistration{provider: pf, config: pcf}
	defaultProviderRegistry.register(name, registration)
}

func (r *providerRegistry) createProvider(name string, config interface{}) (Provider, error) {
	r.mu.RLock()
	var p *providerRegistration
	if r.registrations != nil {
		p = r.registrations[name]
	}
	r.mu.RUnlock()
	if p == nil {
		return nil, ErrProviderNotFound(name)
	}
	return p.provider(config)
}

func (r *providerRegistry) createConfig(name string) (interface{}, error) {
	r.mu.RLock()
	var p *providerRegistration
	if r.registrations != nil {
		p = r.registrations[name]
	}
	r.mu.RUnlock()
	if p == nil {
		return nil, ErrProviderNotFound(name)
	}
	return p.config(), nil
}
