package jsotto_test

import (
	"testing"

	"context"

	"github.com/influx6/faux/tests"
	"github.com/influx6/geckodataset/dataset/config"
	"github.com/influx6/geckodataset/dataset/procs/jsotto"
)

func TestJSOtto(t *testing.T) {
	jt, err := jsotto.New(config.JSOttoConf{
		Main:   "./fixtures/main.js",
		Target: "ParseRecord",
	})

	if err != nil {
		tests.FailedWithError(err, "Should have successfully created JSOtto instance")
	}
	tests.Passed("Should have successfully created JSOtto instance")

	res, err := jt.Transform(context.Background(), map[string]interface{}{
		"age":  20,
		"name": "Alex Woldart",
	})

	if err != nil {
		tests.FailedWithError(err, "Should have successfully transformed data")
	}
	tests.Passed("Should have successfully transformed data")

	if len(res) == 0 {
		tests.Failed("Should have received atleast 1 record")
	}
	tests.Passed("Should have received atleast 1 record")

	total, ok := res[0]["total"].(float64)
	if !ok {
		tests.Failed("Should have found 'total' key in result")
	}
	tests.Passed("Should have found 'total' key in result")

	if total != 1 {
		tests.Failed("Should have matched total to 1")
	}
	tests.Passed("Should have matched total to 1")

}
