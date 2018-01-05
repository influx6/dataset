package config

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
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
	DefaultInterval = time.Second * 60
)

// DriverConfig embodies the configuration used for defining user driver processor.
type DriverConfig struct {
	// JS indicates the configuration values for the JSOtto procs.
	JS *JSOttoConf `toml:"js"`

	// Binary indicates the configuration values to be used for the BinaryRunc procs.
	Binary *BinaryConf `toml:"binary"`

	// Driver value indicates which proc is to be used for processing: js or binary.
	Driver string `toml:"driver"`
}

// Validate returns an error if the config is invalid.
func (dc *DriverConfig) Validate() error {
	if dc.Binary == nil && dc.JS == nil {
		return errors.New("ProcConfig.JS or ProcConfig.Binary must either be provided")
	}

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
	APIKey string `toml:"api_key"`

	// Pull, process and update record at giving intervals. (Optional)
	Interval string `toml:"interval"`

	// PullBatch indicates total records expected by proc to be processed.
	PullBatch int `toml:"pull_batch"`

	// PushBatch indicates total records to be pushed per call to the upstream API.
	PushBatch int `toml:"pull_batch"`

	// RunInterval gets the interval value provided through the `Interval` field or
	// is set to DefaultInterval.
	RunInterval time.Duration `toml:"-"`
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
		return errors.New("APIKey is required")
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
	Name     string `toml:"name"`
	Type     string `toml:"type"`
	Currency string `toml:"currency"`
	Optional bool   `toml:"optional"`
}

// DatasetConfig embodies the configuration data used to define the dataset to
// be used and corresponding dataset field values to be used to create dataset.
type DatasetConfig struct {
	// Dataset indicates the dataset to be used for saving processed results.
	Dataset string `toml:"dataset"`

	// Fields indicates the fields defining the dataset which is expected to be used
	// for storing the processed records.
	Fields []FieldType `toml:"fields"`
}

// Validate returns an error if the config is invalid.
func (dc *DatasetConfig) Validate() error {
	if dc.Dataset == "" {
		return errors.New("Dataset name is required")
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
