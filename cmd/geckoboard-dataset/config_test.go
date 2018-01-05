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

				core := list.JSONFiles[0]
				if core.Source != "./fixtures/sales/user_sales.json" {
					tests.Failed("Should have matched provided source")
				}
				tests.Passed("Should have matched provided source")
			},
		},
		{
			Config: `
interval= "60s"
pull_batch = 100
push_batch = 100
api_key = "your_api_key"

[[datasets]]
driver = "mongodb"

[datasets.conf]
dataset = "user_sales_freq"
source = "user_sales_collection"
dest = "user_sales_metrics" # optional, we want to save transformed records here

[datasets.conf.db]
authdb = "admin"
db = "machines_sales"
user = "tobi_mach"
password = "xxxxxxxxxxxx"
host = "db.mongo.com:4500"

[datasets.conf.binary]
binary = "echo"

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
				if len(list.Mongo) == 0 {
					tests.Failed("Should have passed configuration for json file")
				}
				tests.Passed("Should have passed configuration for json file")

				core := list.Mongo[0]
				if core.Binary == nil {
					tests.Failed("Should have received binary config")
				}
				tests.Passed("Should have received binary config")

				if core.Binary.Binary != "/bin/echo" {
					tests.Failed("Should have received binary command as 'echo'")
				}
				tests.Passed("Should have received binary command as 'echo'")

				if core.DB.Host != "db.mongo.com:4500" {
					tests.Failed("Should have gotten mongodb host")
				}
				tests.Passed("Should have gotten mongodb host")

				if core.DB.AuthDB != "admin" {
					tests.Failed("Should have gotten core.db.authdb")
				}
				tests.Passed("Should have gotten core.db.authdb")

				if core.DB.DB != "machines_sales" {
					tests.Failed("Should have gotten core.db.db")
				}
				tests.Passed("Should have gotten core.db.db")

				if core.DB.User != "tobi_mach" {
					tests.Failed("Should have gotten core.db.user")
				}
				tests.Passed("Should have gotten core.db.user")
			},
		},
		{
			Config: `
interval= "60s"
pull_batch = 100
push_batch = 100
api_key = "your_api_key"

[[datasets]]
driver = "json-file"

[datasets.conf]
dataset = "user_sales_freq"
source = "./fixtures/sales/user_sales.json"

[datasets.conf.binary]
binary = "echo"

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

				core := list.JSONFiles[0]
				if core.Binary == nil {
					tests.Failed("Should have received binary config")
				}
				tests.Passed("Should have received binary config")

				if core.Binary.Binary != "/bin/echo" {
					tests.Failed("Should have received binary command as 'echo'")
				}
				tests.Passed("Should have received binary command as 'echo'")
			},
		},
		{
			Config: `
interval= "60s"
pull_batch = 100
push_batch = 100
api_key = "your_api_key"

[[datasets]]
driver = "json-dir"

[datasets.conf]
dataset = "user_sales_freq"
source_dir = "./fixtures/sales"

[datasets.conf.binary]
binary = "echo"

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
				if len(list.JSONDirs) == 0 {
					tests.Failed("Should have passed configuration for json file")
				}
				tests.Passed("Should have passed configuration for json file")

				core := list.JSONDirs[0]
				if core.Binary == nil {
					tests.Failed("Should have received binary config")
				}
				tests.Passed("Should have received binary config")

				if core.Binary.Binary != "/bin/echo" {
					tests.Failed("Should have received binary command as 'echo'")
				}
				tests.Passed("Should have received binary command as 'echo'")

				if core.SourceDir != "./fixtures/sales" {
					tests.Failed("Should have directory pointing to sales")
				}
				tests.Passed("Should have directory pointing to sales")
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
