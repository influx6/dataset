package binary_test

import (
	"os"
	"testing"

	"context"

	"github.com/influx6/faux/metrics"
	"github.com/influx6/faux/metrics/custom"
	"github.com/influx6/faux/tests"
	"github.com/influx6/geckodataset/dataset/config"
	"github.com/influx6/geckodataset/dataset/procs/binary"
)

func TestBinaryRun(t *testing.T) {
	events := metrics.New()
	if testing.Verbose() {
		events = metrics.New(custom.StackDisplay(os.Stdout))
	}

	binrunc := binary.New(config.BinaryConf{
		Binary: "./fixtures/bin/grun",
	}, events)

	res, err := binrunc.Transform(context.Background(), map[string]interface{}{
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

func TestBinaryRunWithCommand(t *testing.T) {
	events := metrics.New()
	if testing.Verbose() {
		events = metrics.New(custom.StackDisplay(os.Stdout))
	}

	binrunc := binary.New(config.BinaryConf{
		Binary:  "./fixtures/bin/trun",
		Command: "Transform",
	}, events)

	res, err := binrunc.Transform(context.Background(), map[string]interface{}{
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
