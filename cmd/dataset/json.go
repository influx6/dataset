package main

import (
	"errors"
	"os"
	"strings"
	"time"

	"context"

	"github.com/influx6/dataset/dataset"
	"github.com/influx6/dataset/dataset/config"
	"github.com/influx6/dataset/dataset/procs/binary"
	"github.com/influx6/dataset/dataset/procs/jsotto"
	"github.com/influx6/dataset/dataset/pullers/jsonfiles"
	"github.com/influx6/dataset/dataset/pushers"
	"github.com/influx6/faux/metrics"
)

func runDataset(ctx context.Context, conf jsonDataset, base config.ProcConfig) error {
	geckoboard, err := pushers.NewGeckoboardPusher(conf.Dataset)
	if err != nil {
		return err
	}

	stream, err := jsonfiles.NewJSONStream(conf.Source)
	if err != nil {
		return err
	}

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

	var pushers dataset.DataPushers
	pushers = append(pushers, geckoboard)

	controller := dataset.Dataset{
		Pull:    &stream,
		Pushers: pushers,
		Proc:    transformer,
	}

	for {
		// Seek new batch for processing.
		if err := controller.Do(ctx, base.PullBatch, base.PushBatch); err != nil {
			if err == dataset.ErrNoMore {
				return nil
			}

			return err
		}

		// Sleep for giving duration after last run of pull-process-push routine.
		time.Sleep(base.RunInterval)
	}
}

// jsonDataset defines json dataset requests for
// specific file.
type jsonDataset struct {
	config.DriverConfig
	Source  string               `toml:"source"`
	Dataset config.DatasetConfig `toml:"datasets"`
}

// Validate returns an error if the config is invalid.
func (c *jsonDataset) Validate() error {
	if err := c.DriverConfig.Validate(); err != nil {
		return err
	}

	if c.Source != "" {
		return errors.New("config.Source must be provided")
	}

	stat, err := os.Stat(c.Source)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return errors.New("config.Source must be a file")
	}

	if err := c.Dataset.Validate(); err != nil {
		return err
	}

	return nil
}
