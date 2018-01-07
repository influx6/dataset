package main

import (
	"context"
	"errors"

	"time"

	"github.com/influx6/faux/db/mongo"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/geckodataset/dataset"
	"github.com/influx6/geckodataset/dataset/config"
	"github.com/influx6/geckodataset/dataset/procs/binary"
	"github.com/influx6/geckodataset/dataset/procs/jsotto"
	"github.com/influx6/geckodataset/dataset/pushers"
)

func runMGODataset(ctx context.Context, set config.DatasetConfig, ds mgoDataset, conf config.ProcConfig) error {
	if ds.JS == nil && ds.Binary == nil {
		return errors.New("JS or Binary configuration required")
	}

	geckoboard, err := pushers.NewGeckoboardPusher(conf.APIKey, set)
	if err != nil {
		return err
	}

	mdb := mongo.NewMongoDB(ds.DB)

	var transformer dataset.Proc

	if ds.Binary != nil {
		transformer = binary.New(*ds.Binary, metrics.New())
	}

	if ds.JS != nil {
		jso, err := jsotto.New(*ds.JS)
		if err != nil {
			return err
		}

		transformer = jso
	}

	puller := new(mongo.MongoPull)
	puller.Src = mdb
	puller.Collection = ds.Source

	var pushers dataset.DataPushers
	pushers = append(pushers, geckoboard)

	if ds.Destination != "" {
		var mgopusher mongo.MongoPush
		mgopusher.Src = mdb
		mgopusher.Collection = ds.Destination
		pushers = append(pushers, mgopusher)
	}

	controller := dataset.Dataset{
		Pull:    puller,
		Pushers: pushers,
		Proc:    transformer,
	}

	for {
		// Seek new batch for processing.
		if err := controller.Do(ctx, conf.PullBatch, conf.PushBatch); err != nil {
			if err == dataset.ErrNoMore {
				return nil
			}

			return err
		}

		// Sleep for giving duration after last run of pull-process-push routine.
		time.Sleep(conf.RunInterval)
	}
}

// mgoDataset defines json dataset requests for
// specific file.
type mgoDataset struct {
	config.DriverConfig
	config.DatasetConfig

	Destination string       `toml:"dest" json:"dest"`
	Source      string       `toml:"source" json:"source"`
	DB          mongo.Config `toml:"db" json:"db"`
}

// Validate returns an error if the config is invalid.
func (c *mgoDataset) Validate() error {
	if err := c.DriverConfig.Validate(); err != nil {
		return err
	}

	if err := c.DatasetConfig.Validate(); err != nil {
		return err
	}

	if c.Source == "" {
		return errors.New("mongo.Source is required")
	}

	if err := c.DB.Validate(); err != nil {
		return err
	}

	return nil
}
