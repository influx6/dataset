package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"context"

	"github.com/influx6/faux/metrics"
	"github.com/influx6/geckodataset/dataset"
	"github.com/influx6/geckodataset/dataset/config"
	"github.com/influx6/geckodataset/dataset/procs/binary"
	"github.com/influx6/geckodataset/dataset/procs/jsotto"
	"github.com/influx6/geckodataset/dataset/pullers/jsonfiles"
	"github.com/influx6/geckodataset/dataset/pushers"
)

func runJSONDataset(ctx context.Context, set config.DatasetConfig, conf jsonDataset, base config.ProcConfig) error {
	if conf.JS == nil && conf.Binary == nil {
		return errors.New("JS or Binary configuration required")
	}

	geckoboard, err := pushers.NewGeckoboardPusher(base.APIKey, set)
	if err != nil {
		return err
	}

	stream, err := jsonfiles.NewJSONStream(conf.Source)
	if err != nil {
		return err
	}

	var transformer dataset.Proc
	if conf.Binary != nil {
		transformer = binary.New(*conf.Binary, metrics.New())
	}

	if conf.JS != nil {
		jso, err := jsotto.New(*conf.JS)
		if err != nil {
			return err
		}

		transformer = jso
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
	config.DatasetConfig

	Source string `toml:"source"`
}

// Validate returns an error if the config is invalid.
func (c *jsonDataset) Validate() error {
	if err := c.DriverConfig.Validate(); err != nil {
		return err
	}

	if err := c.DatasetConfig.Validate(); err != nil {
		return err
	}

	if c.Source == "" {
		return errors.New("config.Source must be provided")
	}

	stat, err := os.Stat(c.Source)
	if err != nil {
		return fmt.Errorf("config.Source %+q failed to be found", c.Source)
	}

	if stat.IsDir() {
		return errors.New("config.Source must be a file")
	}

	return nil
}
