package mgo_dataset

import "github.com/influx6/faux/flags"

func main() {
	flags.Run("mgodataset", flags.Command{
		Name:      "push",
		ShortDesc: "push data from a mongodb to the geckoboard API",
		Desc: `MgoDataset provides a CLI tooling to allow pushing data from a
mongodb collection to user's Geckobaord API account`,
		Flags: []flags.Flag{},
		Action: func(context flags.Context) error {

			return nil
		},
	})
}
