package main

import (
	"github.com/BurntSushi/toml"
	"github.com/influx6/dataset/dataset/config"
	"github.com/influx6/faux/db/mongo"
	"github.com/influx6/faux/flags"
)

// Config embodies the configuration expected to be loaded
// by user for processing a collection which would then be
// saved to the Geckoboard API.
type Config struct {
	config.DatasetConfig
	Dest   *mongo.Config `toml:"dest"`
	Source mongo.Config  `toml:"source"`
}

// Load attempts to use toml to decode file content into Config instance.
func (c *Config) Load(targetFile string) error {
	if _, err := toml.DecodeFile(targetFile, c); err != nil {
		return err
	}

	return c.Validate()
}

// Validate returns an error if the config is invalid.
func (c *Config) Validate() error {
	if err := c.Source.Validate(); err != nil {
		return err
	}

	if c.Dest != nil {
		if err := c.Dest.Validate(); err != nil {
			// if the Collection is not set and we are still not
			// empty then we have a configuration error.
			if c.Dest.Collection == "" && !c.Dest.Empty() {
				return err
			}

			// if the destination collection is set, then we are
			// properly dealing we new collection to house processed
			// result, but should use existing Source credentials.
			// So we copy c.Source then change collection
			if c.Dest.Collection != "" {
				newDest := c.Source.CloneWithCollection(c.Dest.Collection)
				c.Dest = &newDest
			}
		}
	}

	return c.DatasetConfig.Validate()
}

func main() {
	flags.Run("mgo-dataset", flags.Command{
		Name:      "push",
		ShortDesc: "push data from a mongodb to the geckoboard API",
		Desc: `MgoDataset provides a CLI tooling to allow pushing data from a
mongodb collection to user's Geckobaord API account`,
		Flags: []flags.Flag{},
		Action: func(context flags.Context) error {

			return nil
		},
	})
}
