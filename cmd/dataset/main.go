package main

import (
	"github.com/influx6/faux/flags"
)

func main() {
	flags.Run("dataset", flags.Command{
		Name:      "mongo",
		Action:    mgoAction,
		ShortDesc: "Push data from a mongodb to the geckoboard API",
		Desc:      `mongo allows pushing data from a MongoDB database collection after getting processed using the dataset system then pushed to user's Geckoboard API account`,
		Flags: []flags.Flag{
			&flags.StringFlag{
				Name:    "config",
				Default: "config.toml",
			},
		},
	}, flags.Command{
		Name:      "json-file",
		Action:    jsonAction,
		ShortDesc: "Push data from a file to the geckoboard API",
		Desc:      `json allows pushing data from a json file which are processed using the dataset system then pushed to user's Geckoboard API account`,
		Flags: []flags.Flag{
			&flags.StringFlag{
				Name:    "config",
				Default: "config.toml",
			},
		},
	}, flags.Command{
		Name:      "json-dir",
		Action:    jsonDirAction,
		ShortDesc: "Push data from a directory of json files to the geckoboard API",
		Desc:      `json allows pushing data from a collection of json files which are processed using the dataset system then pushed to user's Geckoboard API account`,
		Flags: []flags.Flag{
			&flags.StringFlag{
				Name:    "config",
				Default: "config.toml",
			},
		},
	})
}
