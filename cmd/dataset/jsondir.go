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

func runDirDataset(ctx context.Context, conf jsonDirDataset, base config.ProcConfig) error {
	geckoboard, err := pushers.NewGeckoboardPusher(conf.Dataset)
	if err != nil {
		return err
	}

	stream, err := jsonfiles.New(conf.SourceDir, conf.Deep)
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
	SourceDir string               `toml:"source"`
	Deep      bool                 `toml:"deep"`
	Dataset   config.DatasetConfig `toml:"datasets"`
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

	if err := c.Dataset.Validate(); err != nil {
		return err
	}

	return nil
}
