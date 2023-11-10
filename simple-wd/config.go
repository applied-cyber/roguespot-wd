package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const configFile = "config.yaml" // Default config file name

type Configuration struct {
	InterfaceName       string `yaml:"interfaceName"`
	ScanIntervalSeconds uint64 `yaml:"scanInterval"`
}

func defaultConfig() *Configuration {
	return &Configuration{
		InterfaceName:       "wlan0",
		ScanIntervalSeconds: 60,
	}
}

func NewConfig() *Configuration {
	newConfig := defaultConfig()

	data, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Failed to read config file '%s': %s\nUsing default configuration values\n",
			configFile,
			err,
		)
		return newConfig
	}

	// yaml.Unmarshal applies the YAML config to the config object
	if err = yaml.Unmarshal(data, &newConfig); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Failed to parse YAML in config file '%s': %s\nUsing default configuration values\n",
			configFile,
			err,
		)
	}

	return newConfig
}
