package main

import (
	"testing"

	"context"

	"github.com/influx6/faux/tests"
)

func TestLoadConfig(t *testing.T) {
	configs := []struct {
		Config   string
		DoError  func(error)
		DoAction func(list datasetList)
	}{
		{
			Config: `interval= "60s"
pull_batch = 100
push_batch = 100
api_key = "your_api_key"

[[datasets]]
driver = "json-file"

[datasets.conf.js]
target = "transformDocument"
main = "./fixtures/transforms/js/user_sales.js"
libraries = ["./fixtures/transforms/js/support/types.js"]

[datasets.conf]
dataset = "user_sales_freq"
source = "./fixtures/sales/user_sales.json"

[[datasets.conf.fields]]
name = "user"
type = "string"

[[datasets.conf.fields]]
name = "scores"
type = "number"

`,
			DoError: func(err error) {
				if err != nil {
					tests.FailedWithError(err, "Should have successfully loaded config")
				}
				tests.Passed("Should have successfully loaded config")
			},
			DoAction: func(list datasetList) {
				if len(list.JSONFiles) == 0 {
					tests.Failed("Should have passed configuration for json file")
				}
				tests.Passed("Should have passed configuration for json file")
			},
		},
	}

	for _, t := range configs {
		res, err := loadConfig(context.Background(), t.Config)
		if t.DoError != nil {
			t.DoError(err)
		}
		if t.DoAction != nil {
			t.DoAction(res)
		}
	}
}
