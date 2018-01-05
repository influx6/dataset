package main

import (
	"github.com/influx6/faux/flags"
)

func main() {
	flags.Run("dataset", flags.Command{
		Name:      "push",
		ShortDesc: "Push data from a mongodb to the geckoboard API",
		Desc:      `Push takes provided configuration which uses to retrieve, process and push new data to user's dataset repositories on the Geckoboard API.`,
		Action: func(context flags.Context) error {
			configFile, _ := context.GetString("config")
			return pushDatasets(context, configFile)
		},
		Flags: []flags.Flag{
			&flags.StringFlag{
				Name:    "config",
				Default: "config.toml",
			},
		},
	})
}
