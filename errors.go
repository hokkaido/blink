package blink

import (
	"fmt"
)

func ErrProviderNotFound(providerName string) error {
	return fmt.Errorf("A provider with name %q was not found", providerName)
}

func ErrProviderAlreadyRegistered(providerName string) error {
	return fmt.Errorf("A provider with name %q is already registered", providerName)
}

func ErrLayerNotFound(providerName string) error {
	return fmt.Errorf("A layer with name %q was not found", providerName)
}

func ErrLayerAlreadyRegistered(providerName string) error {
	return fmt.Errorf("A layer with name %q is already registered", providerName)
}

func ErrLayerConfigParseError(layerName string, innerError error) error {
	return fmt.Errorf("An error occured while parsing the configuration section of layer %q: %q", layerName, innerError)
}

func ErrServerConfigParseError(configFile string, innerError error) error {
	return fmt.Errorf("An error occured while parsing the configuration file %q", configFile)
}
