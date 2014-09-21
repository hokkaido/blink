package blink

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ServerConfig struct {
	Cache  *CacheConfig
	Layers map[string]*LayerConfig
}

type CacheConfig struct {

	// Name of the Cache to use
	Name string
}

type LayerConfig struct {
	Provider   *json.RawMessage
	Projection string
}

type providerNameHolder struct {
	Name string
}

// Initiates the Load Process
func (b *Blink) parseConfig(configFile string) (*ServerConfig, error) {

	var serverConfig *ServerConfig

	f, err := ioutil.ReadFile(configFile)

	if err != nil {
		return nil, ErrServerConfigParseError(configFile, err)
	}

	err = json.Unmarshal(f, &serverConfig)

	if err != nil {
		return nil, ErrServerConfigParseError(configFile, err)
	}

	return serverConfig, nil

}

func (b *Blink) loadConfig(serverConfig *ServerConfig) error {
	log.Println("Initating load sequence")
	for layerName, layerConfig := range serverConfig.Layers {
		log.Printf("Initiating load procedure for layer %q", layerName)
		var providerName providerNameHolder
		raw, _ := layerConfig.Provider.MarshalJSON()

		err := json.Unmarshal(raw, &providerName)
		if err != nil {
			return ErrLayerConfigParseError(layerName, err)
		}
		log.Printf("Trying to create provider config for provider with name %q", providerName.Name)
		providerConfig, err := defaultProviderRegistry.createConfig(providerName.Name)
		if err != nil {
			return ErrLayerConfigParseError(layerName, err)
		}

		err = json.Unmarshal(raw, providerConfig)
		if err != nil {
			return ErrLayerConfigParseError(layerName, err)
		}
		log.Println(providerConfig)
		log.Printf("Trying to create provider with name %q", providerName.Name)
		provider, err := defaultProviderRegistry.createProvider(providerName.Name, providerConfig)
		if err != nil {
			return ErrLayerConfigParseError(layerName, err)
		}

		log.Printf("Trying to create layer %q with provider %q", layerName, providerName.Name)

		var layer = &Layer{Name: layerName, Provider: provider}
		b.RegisterLayer(layer)

	}
	return nil
}
