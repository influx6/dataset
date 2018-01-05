package main

import (
	"fmt"

	"strings"

	"time"

	"github.com/BurntSushi/toml"
	"github.com/influx6/dataset/dataset"
	"github.com/influx6/dataset/dataset/config"
	"github.com/influx6/dataset/dataset/procs/binary"
	"github.com/influx6/dataset/dataset/procs/jsotto"
	"github.com/influx6/dataset/dataset/pushers"
	"github.com/influx6/faux/db/mongo"
	"github.com/influx6/faux/flags"
	"github.com/influx6/faux/metrics"
)

func mgoAction(context flags.Context) error {
	configFile, _ := context.GetString("config")

	var conf mgoConfig
	if err := conf.Load(configFile); err != nil {
		return err
	}

	geckoboard, err := pushers.NewGeckoboardPusher(conf.Dataset)
	if err != nil {
		return err
	}

	mdb := mongo.NewMongoDB(conf.Source)

	var transformer dataset.Proc
	switch strings.ToLower(conf.Driver) {
	case "js", "jsotto":
		jso, err := jsotto.New(*conf.JS)
		if err != nil {
			return err
		}

		transformer = jso
	case "binary":
		transformer = binary.New(*conf.Binary, metrics.New())
	}

	puller := new(mongo.MongoPull)
	puller.Src = mdb

	var pushers dataset.DataPushers
	pushers = append(pushers, geckoboard)

	controller := dataset.Dataset{
		Pull:    puller,
		Pushers: pushers,
		Proc:    transformer,
	}

	for {
		// Seek new batch for processing.
		if err := controller.Do(context, conf.PullBatch, conf.PushBatch); err != nil {
			if err == dataset.ErrNoMore {
				return nil
			}

			return err
		}

		// Sleep for giving duration after last run of pull-process-push routine.
		time.Sleep(conf.RunInterval)
	}

	return nil
}

// mgoConfig embodies the configuration expected to be loaded
// by user for processing a collection which would then be
// saved to the Geckoboard API.
type mgoConfig struct {
	config.ProcConfig
	Dest    *mongo.Config        `toml:"dest"`
	Source  mongo.Config         `toml:"source"`
	Dataset config.DatasetConfig `toml:"datasets"`
}

// Load attempts to use toml to decode file content into Config instance.
func (c *mgoConfig) Load(targetFile string) error {
	if _, err := toml.DecodeFile(targetFile, c); err != nil {
		return err
	}

	return c.Validate()
}

// Validate returns an error if the config is invalid.
func (c *mgoConfig) Validate() error {
	if err := c.ProcConfig.Validate(); err != nil {
		return err
	}

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

	if err := c.Dataset.Validate(); err != nil {
		return fmt.Errorf("dataset %+q: %+s", ds.Dataset, err.Error())
	}

	return nil
}
