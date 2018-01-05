package main

import (
	"github.com/influx6/faux/flags"
)

func main() {
	flags.Run("dataset", flags.Command{
		Name:      "push",
		Action:    mgoAction,
		ShortDesc: "Push data from a mongodb to the geckoboard API",
		Desc:      `mongo allows pushing data from a MongoDB database collection after getting processed using the dataset system then pushed to user's Geckoboard API account`,
		Flags: []flags.Flag{
			&flags.StringFlag{
				Name:    "config",
				Default: "config.toml",
			},
		},
	})
}
