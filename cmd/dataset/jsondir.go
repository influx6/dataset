package main

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/influx6/dataset/dataset"
	"github.com/influx6/dataset/dataset/config"
	"github.com/influx6/dataset/dataset/procs/binary"
	"github.com/influx6/dataset/dataset/procs/jsotto"
	"github.com/influx6/dataset/dataset/pullers/jsonfiles"
	"github.com/influx6/dataset/dataset/pushers"
	"github.com/influx6/faux/flags"
	"github.com/influx6/faux/metrics"
)

func jsonDirAction(context flags.Context) error {
	configFile, _ := context.GetString("config")

	var conf jsonDirConfig
	if err := conf.Load(configFile); err != nil {
		return err
	}

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

// jsonDirConfig embodies the configuration expected to be loaded
// by user for processing a collection which would then be
// saved to the Geckoboard API.
type jsonDirConfig struct {
	config.ProcConfig
	Deep      bool                 `toml"deep"`
	SourceDir string               `toml:"source_dir"`
	Dataset   config.DatasetConfig `toml:"datasets"`
}

// Load attempts to use toml to decode file content into Config instance.
func (c *jsonDirConfig) Load(targetFile string) error {
	if _, err := toml.DecodeFile(targetFile, c); err != nil {
		return err
	}

	return c.Validate()
}

// Validate returns an error if the config is invalid.
func (c *jsonDirConfig) Validate() error {
	if err := c.ProcConfig.Validate(); err != nil {
		return err
	}

	if c.SourceDir != "" {
		return errors.New("config.SourceDir must be provided")
	}

	stat, err := os.Stat(c.SourceDir)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return errors.New("config.SourceDir must be a file")
	}

	if err := c.Dataset.Validate(); err != nil {
		return err
	}

	return nil
}
