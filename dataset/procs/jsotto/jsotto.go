package jsotto

import (
	"context"
	"errors"
	"io/ioutil"

	"github.com/hashicorp/packer/common/json"
	"github.com/influx6/geckodataset/dataset/config"
	"github.com/robertkrimen/otto"
)

// JSOtto implements the the Procs interface and implements processing of
// map records received from a external input source. It uses a loaded
// javascript file and associated libraries with a target function which
// when called accepts a json and processes to return a json of the expected
// output.
// This allows non-go developers, to quickly write transforms in JS which
// transforms data easily.
type JSOtto struct {
	Fn   otto.Value
	VM   *otto.Otto
	Conf config.JSOttoConf
}

// New returns a new instance of JSOtto which implements the Procs interface.
func New(conf config.JSOttoConf) (JSOtto, error) {
	vm := otto.New()

	// Attempt to load all libraries first into vm,
	// return error if error occured.
	for _, library := range conf.Libraries {
		libdata, err := ioutil.ReadFile(library)
		if err != nil {
			return JSOtto{}, err
		}

		_, err = vm.Run(libdata)
		if err != nil {
			return JSOtto{}, err
		}
	}

	// Attempt to load main library into vm,
	// return error if it occurs also.
	maindata, err := ioutil.ReadFile(conf.Main)
	if err != nil {
		return JSOtto{}, err
	}

	_, err = vm.Run(maindata)
	if err != nil {
		return JSOtto{}, err
	}

	fn, err := vm.Get(conf.Target)
	if err != nil {
		return JSOtto{}, err
	}

	if !fn.IsFunction() {
		return JSOtto{}, errors.New("JSOttoConf.Target must be a function")
	}

	return JSOtto{
		Fn:   fn,
		VM:   vm,
		Conf: conf,
	}, nil
}

// Transforms takes incoming records which it transforms into json then calls appropriate
func (jso JSOtto) Transform(ctx context.Context, records ...map[string]interface{}) ([]map[string]interface{}, error) {
	jsonr, err := jso.VM.Get("JSON")
	if err != nil {
		return nil, err
	}

	recJSON, err := jsonr.Object().Call("stringify", records)
	if err != nil {
		return nil, err
	}

	resJSON, err := jso.Fn.Call(jso.Fn, recJSON)
	if err != nil {
		return nil, err
	}

	resJSONExported, err := resJSON.Export()
	if err != nil {
		return nil, err
	}

	if resJSONEx, ok := resJSONExported.(string); ok {
		var rex []map[string]interface{}
		if err := json.Unmarshal([]byte(resJSONEx), &rex); err != nil {
			return nil, err
		}

		return rex, nil
	}

	if resJSONEx, ok := resJSONExported.([]map[string]interface{}); ok {
		return resJSONEx, nil
	}

	return nil, errors.New("invalid type received")
}
