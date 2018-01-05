package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/influx6/geckodataset/dataset/config"
)

type datasetConfig struct {
	config.DatasetConfig
	Driver string                 `toml:"driver"`
	Conf   map[string]interface{} `toml:"conf"`
}

type datasetList struct {
	Config    pushConfig
	Mongo     []mgoDataset
	JSONFiles []jsonDataset
	JSONDirs  []jsonDirDataset
}

type pushConfig struct {
	config.ProcConfig
	Datasets []datasetConfig `json:"datasets"`
}

func loadConfig(ctx context.Context, configData string) (datasetList, error) {
	var config pushConfig
	if _, err := toml.Decode(configData, &config); err != nil {
		return datasetList{}, err
	}

	if err := config.ProcConfig.Validate(); err != nil {
		return datasetList{}, err
	}

	var dl datasetList
	dl.Config = config

	var encoded bytes.Buffer
	for _, dataset := range config.Datasets {
		encoded.Reset()
		if err := toml.NewEncoder(&encoded).Encode(dataset.Conf); err != nil {
			return datasetList{}, err
		}

		switch strings.ToLower(dataset.Driver) {
		case "mongodb":
			var mconf mgoDataset
			if _, err := toml.Decode(encoded.String(), &mconf); err != nil {
				return datasetList{}, err
			}

			if err := mconf.Validate(); err != nil {
				return datasetList{}, err
			}

			mconf.DatasetConfig = dataset.DatasetConfig
			dl.Mongo = append(dl.Mongo, mconf)
		case "json-dir":
			var jsondirconf jsonDirDataset
			if _, err := toml.Decode(encoded.String(), &jsondirconf); err != nil {
				return datasetList{}, err
			}

			if err := jsondirconf.Validate(); err != nil {
				return datasetList{}, err
			}

			jsondirconf.DatasetConfig = dataset.DatasetConfig
			dl.JSONDirs = append(dl.JSONDirs, jsondirconf)
		case "json-file":
			var jsonconf jsonDataset
			if _, err := toml.Decode(encoded.String(), &jsonconf); err != nil {
				return datasetList{}, err
			}

			if err := jsonconf.Validate(); err != nil {
				return datasetList{}, err
			}

			jsonconf.DatasetConfig = dataset.DatasetConfig
			dl.JSONFiles = append(dl.JSONFiles, jsonconf)
		}
	}

	return dl, nil
}

func pushDatasets(ctx context.Context, configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	config, err := loadConfig(ctx, string(data))
	if err != nil {
		return err
	}

	for _, conf := range config.Mongo {
		if err := runMGODataset(ctx, conf.DatasetConfig, conf, config.Config.ProcConfig); err != nil {
			return err
		}
	}

	for _, conf := range config.JSONFiles {
		if err := runJSONDataset(ctx, conf.DatasetConfig, conf, config.Config.ProcConfig); err != nil {
			return err
		}
	}

	for _, conf := range config.JSONDirs {
		if err := runJSONDirDataset(ctx, conf.DatasetConfig, conf, config.Config.ProcConfig); err != nil {
			return err
		}
	}

	return nil
}
