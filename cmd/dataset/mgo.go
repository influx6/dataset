package main

import (
	"context"
	"errors"
	"fmt"

	"strings"

	"time"

	"github.com/influx6/dataset/dataset"
	"github.com/influx6/dataset/dataset/config"
	"github.com/influx6/dataset/dataset/procs/binary"
	"github.com/influx6/dataset/dataset/procs/jsotto"
	"github.com/influx6/dataset/dataset/pushers"
	"github.com/influx6/faux/db/mongo"
	"github.com/influx6/faux/metrics"
)

func runMGODataset(ctx context.Context, ds mgoDirDataset, conf config.ProcConfig) error {
	geckoboard, err := pushers.NewGeckoboardPusher(ds.Dataset)
	if err != nil {
		return err
	}

	mdb := mongo.NewMongoDB(ds.SourceConfig)

	var transformer dataset.Proc
	switch strings.ToLower(ds.Driver) {
	case "js", "jsotto":
		jso, err := jsotto.New(*ds.JS)
		if err != nil {
			return err
		}

		transformer = jso
	case "binary":
		transformer = binary.New(*ds.Binary, metrics.New())
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
type mgoDirDataset struct {
	config.DriverConfig
	Source       string
	Destination  string
	SourceConfig mongo.Config         `toml:"source"`
	Dataset      config.DatasetConfig `toml:"datasets"`
}

// Validate returns an error if the config is invalid.
func (c *mgoDirDataset) Validate() error {
	if err := c.Dataset.Validate(); err != nil {
		return fmt.Errorf("dataset %+q: %+s", c.Dataset, err.Error())
	}

	if c.Source == "" {
		return errors.New("mongo.Source is required")
	}

	if err := c.SourceConfig.Validate(); err != nil {
		return err
	}

	return nil
}
