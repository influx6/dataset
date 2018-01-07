package binary

import (
	"bytes"
	"encoding/json"
	"fmt"

	"context"

	"github.com/influx6/faux/exec"
	"github.com/influx6/faux/metrics"
	"github.com/influx6/geckodataset/dataset/config"
)

// BinaryRunc implements dataset.Proc which uses a provided binary file path
// and it's command if provided to use to process record map.
type BinaryRunc struct {
	Config  config.BinaryConf
	Metrics metrics.Metrics
}

// New returns a new instance of
func New(config config.BinaryConf, m metrics.Metrics) *BinaryRunc {
	return &BinaryRunc{
		Metrics: m,
		Config:  config,
	}
}

// Transforms takes incoming records which it transforms into json then calls appropriate
func (br BinaryRunc) Transform(ctx context.Context, records ...map[string]interface{}) ([]map[string]interface{}, error) {
	var input, errs, output bytes.Buffer
	if err := json.NewEncoder(&input).Encode(records); err != nil {
		return nil, err
	}

	command := br.Config.Bin
	if br.Config.Command != "" {
		command = fmt.Sprintf(`%s %s`, br.Config.Bin, br.Config.Command)
	}

	binaryCmd := exec.New(
		exec.Async(),
		exec.Err(&errs),
		exec.Input(&input),
		exec.Output(&output),
		exec.Command(command),
	)

	if err := binaryCmd.Exec(ctx, br.Metrics); err != nil {
		return nil, fmt.Errorf("%+s: %+s", err.Error(), errs.String())
	}

	var res []map[string]interface{}
	if err := json.NewDecoder(&output).Decode(&res); err != nil {
		return nil, fmt.Errorf("%+s: %+s", err.Error(), errs.String())
	}

	return res, nil
}
