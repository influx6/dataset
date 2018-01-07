package main

import (
	"fmt"
	"path/filepath"

	"github.com/influx6/faux/flags"
)

func main() {
	flags.Run("dataset", flags.Command{
		Name:      "push",
		ShortDesc: "Push data from a mongodb to the geckoboard API",
		Desc:      `Push takes provided configuration which uses to retrieve, process and push new data to user's dataset repositories on the Geckoboard API.`,
		Action: func(context flags.Context) error {
			configFile, _ := context.GetString("config")
			switch filepath.Ext(configFile) {
			case ".yaml":
				return pushYAMLDatasets(context, configFile)
			case ".toml":
				return pushTOMLDatasets(context, configFile)
			default:
				return fmt.Errorf("%+q config file extension unknown (support: .yaml, .toml)")
			}
		},
		Flags: []flags.Flag{
			&flags.StringFlag{
				Name:    "config",
				Default: "config.yaml",
				Desc:    "configuration file for processing data into Geckoboard dataset.",
			},
		},
	})
}
