package blink

import (
	"log"
	"sync"
)

// Blink is
type Blink struct {
	config *ServerConfig
	layers map[string]*Layer
	mu     sync.RWMutex
}

type Layer struct {
	Provider Provider
	Name     string
}

// NewBlink allocates and returns a new Blink instance
func New() *Blink {
	return &Blink{layers: make(map[string]*Layer)}
}

// DefaultBlink is the default Blink instance that will be used by Serve
var DefaultBlink = New()

func (b *Blink) Load(configFile string) {
	config, err := b.parseConfig(configFile)
	log.Printf("Configuration %q succesfully parsed", configFile)
	if err != nil {
		log.Printf("An error occured %q", err)
		return

	}
	err = b.loadConfig(config)
	if err != nil {
		log.Printf("An error occured %q", err)
		return

	}

}

func Load(configFile string) {
	DefaultBlink.Load(configFile)
}

// Registers the layer with the given name
// in the default Blink instance.
func RegisterLayer(l *Layer) {
	DefaultBlink.RegisterLayer(l)
}

// Registers a Layer with Blink.
// If a layer with the same name is already registered, Register panics.
func (b *Blink) RegisterLayer(l *Layer) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.layers == nil {
		b.layers = make(map[string]*Layer)
	}
	if _, exists := b.layers[l.Name]; exists {
		panic("A layer with the name " + l.Name + " is already registered")
	}
	b.layers[l.Name] = l
}

// Returns a registered Layer with the given name.
//
// If there is no registered layer with the provided name,
// an is returned error.
func (b *Blink) GetLayer(name string) (*Layer, error) {
	b.mu.RLock()
	var l *Layer
	if b.layers != nil {
		l = b.layers[name]
	}
	b.mu.RUnlock()
	if l == nil {
		return nil, ErrLayerNotFound(name)
	}
	return l, nil
}

// Returns a registered Layer with the given name.
//
// If there is no registered layer with the provided name,
// an is returned error.
func GetLayer(name string) (*Layer, error) {
	return DefaultBlink.GetLayer(name)
}

func NewLayer(name string, provider Provider) *Layer {
	return &Layer{Name: name, Provider: provider}
}

// A convenience function which calls GetTile of the layer's underlying Provider.
func (l *Layer) GetTile(zoom int, x int, y int) ([]byte, error) {
	return l.Provider.GetTile(zoom, x, y)
}
