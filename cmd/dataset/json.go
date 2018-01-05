package main

import (
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/influx6/dataset/dataset/config"
	"github.com/influx6/faux/flags"
)

func jsonAction(context flags.Context) error {
	configFile, _ := context.GetString("config")

	var conf jsonConfig
	if err := conf.Load(configFile); err != nil {
		return err
	}

	return nil
}

// jsonConfig embodies the configuration expected to be loaded
// by user for processing a collection which would then be
// saved to the Geckoboard API.
type jsonConfig struct {
	config.ProcConfig
	Source  string               `toml:"source"`
	Dataset config.DatasetConfig `toml:"datasets"`
}

// Load attempts to use toml to decode file content into Config instance.
func (c *jsonConfig) Load(targetFile string) error {
	if _, err := toml.DecodeFile(targetFile, c); err != nil {
		return err
	}

	return c.Validate()
}

// Validate returns an error if the config is invalid.
func (c *jsonConfig) Validate() error {
	if err := c.ProcConfig.Validate(); err != nil {
		return err
	}

	if c.Source != "" {
		return errors.New("config.Source must be provided")
	}

	if err := c.Dataset.Validate(); err != nil {
		return err
	}

	return nil
}
