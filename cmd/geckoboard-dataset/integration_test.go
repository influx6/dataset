package main

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/influx6/faux/tests"
)

var (
	geckoEnvName = "GECKOBOARD_TEST_KEY"
	testconfig   = `interval: 1s
pull_batch: 3
push_batch: 3
api_key: {{ env "GECKOBOARD_TEST_KEY" }}
datasets:
 - driver: "json-file"
   op: push
   dataset: "user_sales_freq"
   fields:
    - name: user
      type: string
    - name: sales
      type: number
   conf:
    source: "./fixtures/sales/user_sales.json"
    js:
     target: Transform
     main: "./fixtures/transforms/js/user_sales.js"
`
)

func TestJavascriptPushIntegration(t *testing.T) {
	if strings.TrimSpace(os.Getenv(geckoEnvName)) == "" {
		tests.Info("TestCLIRun requires %+q environment variable first", geckoEnvName)
		return
	}

	loadedConfig, err := loadYAMLConfig(context.Background(), testconfig)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully loaded user configuration")
	}
	tests.Passed("Should have successfully loaded user configuration")

	if err := runDatasetConfig(context.Background(), loadedConfig); err != nil {
		tests.FailedWithError(err, "Should have successfully executed configuration for Geckoboard Dataset")
	}
	tests.Passed("Should have successfully executed configuration for Geckoboard Dataset")
}
