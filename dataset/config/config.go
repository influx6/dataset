package config

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/influx6/faux/db/mongo"
)

const (
	// DefaultBatch indicates the total records expected to be
	// sent into proc for processing, which can vise-versa mean
	// the expected records to be produced from output, but this
	// is user dependent and not a hard rule.
	DefaultBatch = 1
)

// Config embodies the configuration expected to be loaded
// by user for processing a collection which would then be
// saved to the Geckoboard API.
type Config struct {
	JS     JSOttoConf    `toml:"js"`
	Dest   *mongo.Config `toml:"dest"`
	Binary BinaryConf    `toml:"binary"`
	Source mongo.Config  `toml:"source"`

	// Pull, process and update record at giving intervals. (Optional)
	Interval string `toml:"interval"`

	// Driver value indicates which proc is to be used for processing: js or binary.
	Driver string `toml:"driver"`

	// Batch indicates total records expected by proc to be processed, default is 1.
	Batch uint16 `toml:"batch"`

	// Dataset indicates the dataset to be used for saving processed results.
	Dataset string `toml:"dataset"`

	// APIKey indicates the user's Geckboard API Key used for authentication of all save requests.
	APIKey string `toml:"api_key"`

	RunInterval time.Duration `toml:"-"`
}

// Load attempts to use toml to decode file content into Config instance.
func (c *Config) Load(targetFile string) error {
	if _, err := toml.DecodeFile(targetFile, c); err != nil {
		return err
	}

	if c.Batch == 0 {
		c.Batch = DefaultBatch
	}

	return c.Validate()
}

// Validate returns an error if the config is invalid.
func (c *Config) Validate() error {
	if err := c.Source.Validate(); err != nil {
		return err
	}

	if c.Dest != nil {
		if err := c.Dest.Validate(); err != nil {
			// if the Collection is not set and we are still not
			// empty then we have a configuration error.
			if c.Dest.Collection == "" && !c.Dest.Empty() {
				return err
			}

			// if the destination collection is set, then we are
			// properly dealing we new collection to house processed
			// result, but should use existing Source credentials.
			// So we copy c.Source then change collection
			if c.Dest.Collection != "" {
				newDest := c.Source.CloneWithCollection(c.Dest.Collection)
				c.Dest = &newDest
			}
		}
	}

	if c.APIKey == "" {
		return errors.New("APIKey is required")
	}

	if c.Dataset == "" {
		return errors.New("Dataset name is required")
	}

	if c.Interval != "" {
		interval, err := time.ParseDuration(c.Interval)
		if err != nil {
			return err
		}
		c.RunInterval = interval
	}

	switch strings.ToLower(c.Driver) {
	case "jsotto", "js":
		return c.JS.Validate()
	case "binary":
		return c.Binary.Validate()
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
