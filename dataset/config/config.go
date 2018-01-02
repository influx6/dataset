package config

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	// DefaultBatch indicates the total records expected to be
	// sent into proc for processing, which can vise-versa mean
	// the expected records to be produced from output, but this
	// is user dependent and not a hard rule.
	DefaultBatch = 100

	// DefaultInterval indicates the default expected time for each
	// requests to be processed before waiting for it's next run.
	DefaultInterval = time.Second * 60
)

type DatasetConfig struct {
	JS     JSOttoConf `toml:"js"`
	Binary BinaryConf `toml:"binary"`

	// Pull, process and update record at giving intervals. (Optional)
	Interval string `toml:"interval"`

	// Driver value indicates which proc is to be used for processing: js or binary.
	Driver string `toml:"driver"`

	// Batch indicates total records expected by proc to be processed, default is 1.
	Batch int `toml:"batch"`

	// Dataset indicates the dataset to be used for saving processed results.
	Dataset string `toml:"dataset"`

	// APIKey indicates the user's Geckboard API Key used for authentication of all save requests.
	APIKey string `toml:"api_key"`

	RunInterval time.Duration `toml:"-"`
}

// Validate returns an error if the config is invalid.
func (dc *DatasetConfig) Validate() error {
	if dc.APIKey == "" {
		return errors.New("APIKey is required")
	}

	if dc.Dataset == "" {
		return errors.New("Dataset name is required")
	}

	if dc.Interval != "" {
		interval, err := time.ParseDuration(dc.Interval)
		if err != nil {
			return err
		}
		dc.RunInterval = interval
	} else {
		dc.RunInterval = DefaultInterval
	}

	if dc.Batch == 0 {
		dc.Batch = DefaultBatch
	}

	switch strings.ToLower(dc.Driver) {
	case "jsotto", "js":
		return dc.JS.Validate()
	case "binary":
		return dc.Binary.Validate()
	}

	return nil
}

// JSOttoConf embodies data used to define the javascript files used for
// providing user processing function for conversion of incoming mongo data
// using the otto javascript vm. https://github.com/robertkrimen/otto.
type JSOttoConf struct {
	Target    string   `toml:"target"`
	Main      string   `toml:"main"`
	Libraries []string `toml:"libraries"`
}

// Validate returns an error if the config is invalid.
func (jsc JSOttoConf) Validate() error {
	if jsc.Target == "" {
		return errors.New("JSOttoConf.Target is required")
	}
	if jsc.Main == "" {
		return errors.New("JSOttoConf.Main is required")
	}
	return nil
}

// BinaryConf embodies data to be used to define the go binary used for processing
// incoming data from the mongo collection.
type BinaryConf struct {
	// Binary path to golang binary for execution, where main expects data coming
	// from stdin with processed data received from stdout.
	Binary string `toml:"binary"`

	// Command name to be used to run against binary if binary is not direct entry
	// point for processor.
	Command string `toml:"command"`
}

// Validate returns an error if the config is invalid.
func (gc *BinaryConf) Validate() error {
	if gc.Binary == "" {
		return errors.New("BinaryConf.Binary is required")
	}

	if filepath.IsAbs(gc.Binary) {
		stat, err := os.Stat(gc.Binary)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			return errors.New("BinaryConf.Binary can't point to a directory")
		}

		return nil
	}

	binaryPath, err := exec.LookPath(gc.Binary)
	if err != nil {
		return errors.New("BinaryConf.Binary target not found in host system: " + err.Error())
	}

	gc.Binary = binaryPath
	return nil
}
