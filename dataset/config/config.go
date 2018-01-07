package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	// DefaultPushBatch indicates the total records expected to be
	// sent to the Pusher, which then delivers to it's corresponding
	// endpoint destination e.g Geckobaord API.
	DefaultPushBatch = 500

	// DefaultPullBatch indicates the total records expected to be
	// pulled and processed by Procs from the source.
	DefaultPullBatch = 500

	// DefaultInterval indicates the default expected time for each
	// requests to be processed before waiting for it's next run.
	DefaultInterval = time.Second * 5
)

// DriverConfig embodies the configuration used for defining user driver processor.
type DriverConfig struct {
	// JS indicates the configuration values for the JSOtto procs.
	JS *JSOttoConf `toml:"js" json:"js"`

	// Binary indicates the configuration values to be used for the BinaryRunc procs.
	Binary *BinaryConf `toml:"binary" json:"binary"`
}

// Validate returns an error if the config is invalid.
func (dc *DriverConfig) Validate() error {
	if dc.JS != nil {
		if err := dc.JS.Validate(); err != nil {
			return err
		}
	}

	if dc.Binary != nil {
		if err := dc.Binary.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// ProcConfig embodies the configuration used for defining user configuration
// for the proc processors who handle conversion of data to datastore records.
type ProcConfig struct {
	// APIKey indicates the user's Geckoboard API Key used for authentication of all save requests.
	APIKey string `toml:"api_key" json:"api_key"`

	// Pull, process and update record at giving intervals. (Optional)
	Interval string `toml:"interval" json:"interval"`

	// PullBatch indicates total records expected by proc to be processed.
	PullBatch int `toml:"pull_batch" json:"pull_batch"`

	// PushBatch indicates total records to be pushed per call to the upstream API.
	PushBatch int `toml:"push_batch" json:"push_batch"`

	// RunInterval gets the interval value provided through the `Interval` field or
	// is set to DefaultInterval.
	RunInterval time.Duration `toml:"-" json:"-"`
}

// Validate returns an error if the config is invalid.
func (dc *ProcConfig) Validate() error {
	if dc.Interval != "" {
		interval, err := time.ParseDuration(dc.Interval)
		if err != nil {
			return err
		}
		dc.RunInterval = interval
	} else {
		dc.RunInterval = DefaultInterval
	}

	if dc.APIKey == "" {
		return errors.New("Config.APIKey is required")
	}

	if dc.PullBatch <= 0 {
		dc.PullBatch = DefaultPullBatch
	}

	if dc.PushBatch <= 0 {
		dc.PushBatch = DefaultPushBatch
	}

	return nil
}

// FieldType embodies field values for defining dataset field types.
type FieldType struct {
	Name     string `toml:"name" json:"name"`
	Type     string `toml:"type" json:"type"`
	Currency string `toml:"currency" json:"currency"`
	Optional bool   `toml:"optional" json:"optional"`
}

// DatasetConfig embodies the configuration data used to define the dataset to
// be used and corresponding dataset field values to be used to create dataset.
type DatasetConfig struct {
	// Operation indicates the type of operation to be performed
	Op string `toml:"op" json:"op"`

	// Dataset indicates the dataset to be used for saving processed results.
	Dataset string `toml:"dataset" json:"dataset"`

	// UniqueBy contains unique values for creating new dataset fields.
	UniqueBy []string `toml:"uniques" json:"unique_by"`

	// DeleteBy contains values used during record updates.
	DeteletBy []string `toml:"delete_by" json:"delete_by"`

	// Fields indicates the fields defining the dataset which is expected to be used
	// for storing the processed records.
	Fields []FieldType `toml:"fields" json:"fields"`
}

// Validate returns an error if the config is invalid.
func (dc *DatasetConfig) Validate() error {
	if dc.Dataset == "" {
		return errors.New("DatasetConfig.Dataset is required")
	}

	if dc.Op == "" {
		dc.Op = "Push"
	}

	if strings.ToLower(dc.Op) != "push" && strings.ToLower(dc.Op) != "update" {
		return fmt.Errorf("DatasetConfig.Op can only be either 'Push' or 'Update' not %q", dc.Op)
	}

	return nil
}

// JSOttoConf embodies data used to define the javascript files used for
// providing user processing function for conversion of incoming mongo data
// using the otto javascript vm. https://github.com/robertkrimen/otto.
type JSOttoConf struct {
	Main      string   `toml:"main" json:"main"`
	Target    string   `toml:"target" json:"target"`
	Libraries []string `toml:"libraries" json:"libraries"`
}

// Validate returns an error if the config is invalid.
func (jsc JSOttoConf) Validate() error {
	if jsc.Target == "" {
		return errors.New("JSOttoConf.Target is required")
	}

	if jsc.Main == "" {
		return errors.New("JSOttoConf.Main is required")
	}

	stat, err := os.Stat(jsc.Main)
	if err != nil {
		return fmt.Errorf("JSOttoConf.Main must exists: %+s", err.Error())
	}

	if stat.IsDir() {
		return errors.New("JSOttoConf.Binary can't point to a directory")
	}

	return nil
}

// BinaryConf embodies data to be used to define the go binary used for processing
// incoming data from the mongo collection.
type BinaryConf struct {
	// Bin path to golang binary for execution, where main expects data coming
	// from stdin with processed data received from stdout.
	Bin string `toml:"bin" json:"bin"`

	// Command name to be used to run against binary if binary is not direct entry
	// point for processor.
	Command string `toml:"command" json:"command"`
}

// Validate returns an error if the config is invalid.
func (gc *BinaryConf) Validate() error {
	if gc.Bin == "" {
		return errors.New("BinaryConf.Binary is required")
	}

	if filepath.IsAbs(gc.Bin) {
		stat, err := os.Stat(gc.Bin)
		if err != nil {
			return fmt.Errorf("BinaryConf.Binary must exists: %+s", err.Error())
		}

		if stat.IsDir() {
			return errors.New("BinaryConf.Binary can't point to a directory")
		}

		return nil
	}

	binaryPath, err := exec.LookPath(gc.Bin)
	if err != nil {
		return errors.New("BinaryConf.Binary target not found in host system: " + err.Error())
	}

	gc.Bin = binaryPath
	return nil
}
