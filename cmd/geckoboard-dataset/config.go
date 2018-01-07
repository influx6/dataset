package main

import (
	"bytes"
	"context"
	"html/template"
	"io/ioutil"
	"strings"

	"os"

	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/ghodss/yaml"
	"github.com/influx6/geckodataset/dataset/config"
)

var (
	tmplFuncs = template.FuncMap{
		"quote": strconv.Quote,
		"env": func(name string) string {
			return strings.TrimSpace(os.Getenv(name))
		},
	}
)

type datasetList struct {
	Config    config.ProcConfig
	Mongo     []mgoDataset
	JSONFiles []jsonDataset
	JSONDirs  []jsonDirDataset
}

type datasetConfig struct {
	config.DatasetConfig
	Driver string                 `toml:"driver" json:"driver"`
	Conf   map[string]interface{} `toml:"conf" json:"conf"`
}

type lConfig struct {
	config.ProcConfig
	Datasets []datasetConfig `toml:"datasets" json:"datasets"`
}

// loadYAMLConfig returns a datasetList which is generated from the provided yaml config
// string returning appropriate config structures.
func loadYAMLConfig(ctx context.Context, configData string) (datasetList, error) {
	var formatted bytes.Buffer
	tml, err := template.New("geckodataset-yaml-config").Funcs(tmplFuncs).Parse(configData)
	if err != nil {
		return datasetList{}, err
	}

	if err := tml.Execute(&formatted, nil); err != nil {
		return datasetList{}, err
	}

	var con lConfig
	if err := yaml.Unmarshal(formatted.Bytes(), &con); err != nil {
		return datasetList{}, err
	}

	if err := con.ProcConfig.Validate(); err != nil {
		return datasetList{}, err
	}

	var dl datasetList
	dl.Config = con.ProcConfig

	for _, dataset := range con.Datasets {
		encoded, err := yaml.Marshal(dataset.Conf)
		if err != nil {
			return datasetList{}, err
		}

		switch strings.ToLower(dataset.Driver) {
		case "mongodb":
			var mconf mgoDataset
			if err := yaml.Unmarshal(encoded, &mconf); err != nil {
				return datasetList{}, err
			}

			mconf.DatasetConfig = dataset.DatasetConfig
			if err := mconf.Validate(); err != nil {
				return datasetList{}, err
			}

			dl.Mongo = append(dl.Mongo, mconf)
		case "json-dir":
			var jsondirconf jsonDirDataset
			if err := yaml.Unmarshal(encoded, &jsondirconf); err != nil {
				return datasetList{}, err
			}

			jsondirconf.DatasetConfig = dataset.DatasetConfig
			if err := jsondirconf.Validate(); err != nil {
				return datasetList{}, err
			}

			dl.JSONDirs = append(dl.JSONDirs, jsondirconf)
		case "json-file":
			var jsonconf jsonDataset
			if err := yaml.Unmarshal(encoded, &jsonconf); err != nil {
				return datasetList{}, err
			}

			jsonconf.DatasetConfig = dataset.DatasetConfig
			if err := jsonconf.Validate(); err != nil {
				return datasetList{}, err
			}

			dl.JSONFiles = append(dl.JSONFiles, jsonconf)
		}
	}

	return dl, nil
}

func pushYAMLDatasets(ctx context.Context, configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	config, err := loadYAMLConfig(ctx, string(data))
	if err != nil {
		return err
	}

	return runDatasetConfig(ctx, config)
}

// loadTOMLConfig returns a datasetList which is generated from the provided toml config
// string returning appropriate config structures.
func loadTOMLConfig(ctx context.Context, configData string) (datasetList, error) {
	var formatted bytes.Buffer
	tml, err := template.New("geckodataset-toml-config").Funcs(tmplFuncs).Parse(configData)
	if err != nil {
		return datasetList{}, err
	}

	if err := tml.Execute(&formatted, nil); err != nil {
		return datasetList{}, err
	}

	var config lConfig
	if _, err := toml.Decode(formatted.String(), &config); err != nil {
		return datasetList{}, err
	}

	if err := config.ProcConfig.Validate(); err != nil {
		return datasetList{}, err
	}

	var dl datasetList
	dl.Config = config.ProcConfig

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

			mconf.DatasetConfig = dataset.DatasetConfig
			if err := mconf.Validate(); err != nil {
				return datasetList{}, err
			}

			dl.Mongo = append(dl.Mongo, mconf)
		case "json-dir":
			var jsondirconf jsonDirDataset
			if _, err := toml.Decode(encoded.String(), &jsondirconf); err != nil {
				return datasetList{}, err
			}

			jsondirconf.DatasetConfig = dataset.DatasetConfig
			if err := jsondirconf.Validate(); err != nil {
				return datasetList{}, err
			}

			dl.JSONDirs = append(dl.JSONDirs, jsondirconf)
		case "json-file":
			var jsonconf jsonDataset
			if _, err := toml.Decode(encoded.String(), &jsonconf); err != nil {
				return datasetList{}, err
			}

			jsonconf.DatasetConfig = dataset.DatasetConfig
			if err := jsonconf.Validate(); err != nil {
				return datasetList{}, err
			}

			dl.JSONFiles = append(dl.JSONFiles, jsonconf)
		}
	}

	return dl, nil
}

func pushTOMLDatasets(ctx context.Context, configFile string) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	config, err := loadTOMLConfig(ctx, string(data))
	if err != nil {
		return err
	}

	return runDatasetConfig(ctx, config)
}

func runDatasetConfig(ctx context.Context, config datasetList) error {
	for _, conf := range config.Mongo {
		if err := runMGODataset(ctx, conf.DatasetConfig, conf, config.Config); err != nil {
			return err
		}
	}

	for _, conf := range config.JSONFiles {
		if err := runJSONDataset(ctx, conf.DatasetConfig, conf, config.Config); err != nil {
			return err
		}
	}

	for _, conf := range config.JSONDirs {
		if err := runJSONDirDataset(ctx, conf.DatasetConfig, conf, config.Config); err != nil {
			return err
		}
	}

	return nil
}
