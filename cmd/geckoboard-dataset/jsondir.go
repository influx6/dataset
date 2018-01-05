package main

import (
	"errors"
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

func runJSONDirDataset(ctx context.Context, set config.DatasetConfig, conf jsonDirDataset, base config.ProcConfig) error {
	if conf.JS == nil && conf.Binary == nil {
		return errors.New("JS or Binary configuration required")
	}

	geckoboard, err := pushers.NewGeckoboardPusher(base.APIKey, set)
	if err != nil {
		return err
	}

	stream, err := jsonfiles.New(conf.SourceDir, conf.Deep)
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
		Pull:    stream,
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

// jsonDirDataset defines json dataset requests for
// specific file.
type jsonDirDataset struct {
	config.DriverConfig
	config.DatasetConfig

	Deep      bool   `toml:"deep"`
	SourceDir string `toml:"source"`
}

// Validate returns an error if the config is invalid.
func (c *jsonDirDataset) Validate() error {
	if err := c.DriverConfig.Validate(); err != nil {
		return err
	}

	if c.SourceDir != "" {
		return errors.New("config.Source must be provided")
	}

	stat, err := os.Stat(c.SourceDir)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return errors.New("config.Source must be a file")
	}

	return nil
}
