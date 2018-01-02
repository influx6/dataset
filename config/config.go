package config

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
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
	Mongo  MongoDBConf `toml:"mongo"`
	Go     GoConf      `toml:"go"`
	JS     JSOttoConf  `toml:"js"`
	Shogun ShogunConf  `toml:"shogun"`

	// Driver value indicates which proc is to be used for processing: js or binary.
	Driver string `toml:"driver"`

	// Batch indicates total records expected by proc to be processed, default is 1.
	Batch uint16 `toml:"batch"`

	// Dataset indicates the dataset to be used for saving processed results.
	Dataset string `toml:"dataset"`

	// APIKey indicates the user's Geckboard API Key used for authentication of all save requests.
	APIKey string `toml:"api_key"`
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
func (c Config) Validate() error {
	if err := c.Mongo.Validate(); err != nil {
		return err
	}

	if c.APIKey == "" {
		return errors.New("APIKey is required")
	}

	if c.Dataset == "" {
		return errors.New("Dataset name is required")
	}

	switch strings.ToLower(c.Driver) {
	case "jsotto", "js":
		return c.JS.Validate()
	case "go":
		return c.Go.Validate()
	case "shogun":
		return c.Shogun.Validate()
	}

	return nil
}

// MongoDBConf embodies the data used to connect to user's mongo connection.
type MongoDBConf struct {
	AuthDB     string `toml:"authdb"`
	DB         string `toml:"db"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	Host       string `toml:"host"`
	Collection string `toml:"collection"`
}

// Validate returns an error if the config is invalid.
func (mgc MongoDBConf) Validate() error {
	if mgc.User == "" {
		return errors.New("MongoDBConf.User is required")
	}
	if mgc.Password == "" {
		return errors.New("MongoDBConf.Password is required")
	}
	if mgc.AuthDB == "" {
		return errors.New("MongoDBConf.AuthDB is required")
	}
	if mgc.Host == "" {
		return errors.New("MongoDBConf.Host is required")
	}
	if mgc.DB == "" {
		return errors.New("MongoDBConf.DB is required")
	}
	return nil
}

// JSOttoConf embodies data used to define the javascript files used for
// providing user processing function for conversion of incoming mongo data
// using the otto javascript vm. https://github.com/robertkrimen/otto.
type JSOttoConf struct {
	Target      string   `toml:"target"`
	Main        string   `toml:"main"`
	JSLibraries []string `toml:"libraries"`
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

// GoConf embodies data to be used to define the go binary used for processing
// incoming data from the mongo collection.
type GoConf struct {
	// Binary Path to golang binary for execution, where main expects data coming
	// from stdin with processed data received from stdout.
	Binary string `toml:"binary"`
}

// Validate returns an error if the config is invalid.
func (gc *GoConf) Validate() error {
	if gc.Binary == "" {
		return errors.New("GoConf.Binary is required")
	}

	if filepath.IsAbs(gc.Binary) {
		stat, err := os.Stat(gc.Binary)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			return errors.New("GoConf.Binary can't point to a directory")
		}

		return nil
	}

	realPath, err := exec.LookPath(gc.Binary)
	if err != nil {
		return errors.New("GoConf.Binary target not found in host system: " + err.Error())
	}

	gc.Binary = realPath
	return nil
}

// ShogunConf embodies data to be used to define the go binary used for processing
// incoming data from the mongo collection.
type ShogunConf struct {
	// Binary Path to golang binary for execution, where main expects data coming
	// from stdin with processed data received from stdout.
	Binary  string `toml:"binary"`
	Command string `toml:"command"`
}

// Validate returns an error if the config is invalid.
func (sc *ShogunConf) Validate() error {
	if sc.Binary == "" {
		return errors.New("GoConf.Binary is required")
	}

	if filepath.IsAbs(sc.Binary) {
		stat, err := os.Stat(sc.Binary)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			return errors.New("GoConf.Binary can't point to a directory")
		}

		return nil
	}

	realPath, err := exec.LookPath(sc.Binary)
	if err != nil {
		return errors.New("GoConf.Binary target not found in host system: " + err.Error())
	}

	sc.Binary = realPath
	return nil
}
