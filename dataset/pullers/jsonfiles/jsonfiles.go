package jsonfiles

import (
	"context"
	"errors"
	"sync"

	"os"

	"bufio"
	"encoding/json"

	"path/filepath"

	"github.com/influx6/geckodataset/dataset"
)

// JSONStream representing in-memory data over processed json data.
// It implements the dataset.Puller interface, and returns records
// as requested in batches.
type JSONStream struct {
	loaded     bool
	targetFile string
	records    []map[string]interface{}
}

// NewJSONStream returns a new instance of JSONStream for giving file.
func NewJSONStream(targetFile string) (JSONStream, error) {
	stat, err := os.Stat(targetFile)
	if err != nil {
		return JSONStream{}, err
	}

	if stat.IsDir() {
		return JSONStream{}, errors.New("only files allowed")
	}

	var js JSONStream
	js.targetFile = targetFile
	return js, nil
}

// load runs the internal processes to lazily load the internal data of giving json file.
func (jns *JSONStream) load(ctx context.Context) error {
	if jns.loaded {
		return nil
	}

	target, err := os.Open(jns.targetFile)
	if err != nil {
		return err
	}

	defer target.Close()

	buffTarget := bufio.NewReader(target)
	if err := json.NewDecoder(buffTarget).Decode(&jns.records); err != nil {
		return err
	}

	jns.loaded = true
	return nil
}

// Pull returns giving set of json records from internal in-memory store which it limits to
// specified batch side.
func (jns *JSONStream) Pull(ctx context.Context, batch int) ([]map[string]interface{}, error) {
	if !jns.loaded {
		if err := jns.load(ctx); err != nil {
			return nil, err
		}
	}

	if batch == 0 {
		return nil, dataset.ErrNoMore
	}

	if len(jns.records) == 0 {
		return nil, dataset.ErrNoMore
	}

	if batch >= len(jns.records) {
		records := jns.records
		jns.records = nil
		return records, nil
	}

	next := jns.records[:batch]
	jns.records = jns.records[batch:]
	return next, nil
}

// JSONStreams embodies the collection of json files loaded from provided directory.
// It creates JSONStream objects which lazy load required data of files during their initial
// call. To pull data.
type JSONStreams struct {
	streams []JSONStream
	ml      sync.Mutex
	current *JSONStream
}

// New returns a new instance of JSONStreams. JSONStreams only generates a file lists of
// files within root if deep is false, else runs into all files with .json prefix.
func New(dir string, deep bool) (*JSONStreams, error) {
	var streams JSONStreams

	if deep {
		if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			if filepath.Ext(path) != ".json" {
				return nil
			}

			stream, err := NewJSONStream(path)
			if err != nil {
				return err
			}

			streams.streams = append(streams.streams, stream)
			return nil
		}); err != nil {
			return nil, err
		}
	} else {
		targetDir, err := os.Open(dir)
		if err != nil {
			return nil, err
		}

		defer targetDir.Close()

		lists, err := targetDir.Readdir(-1)
		if err != nil {
			return nil, err
		}

		for _, item := range lists {
			if item.IsDir() {
				continue
			}

			if filepath.Ext(item.Name()) != ".json" {
				continue
			}

			stream, err := NewJSONStream(filepath.Join(dir, item.Name()))
			if err != nil {
				return &streams, err
			}

			streams.streams = append(streams.streams, stream)
		}
	}

	return &streams, nil
}

// Total returns total records loaded.
func (jns *JSONStreams) Total() int {
	return len(jns.streams)
}

// Pull attempts to load current streams data with batch parameters if found else, walks through
// directory which it loads all fileInfo items, it scans in attempt to load next which if is a valid
// json file and with respect to it's strict flag, will load the content and use this data has
// a means of loading continuous json feed of record values for processing.
func (jns *JSONStreams) Pull(ctx context.Context, batch int) ([]map[string]interface{}, error) {
	jns.ml.Lock()
	defer jns.ml.Unlock()

	if batch == 0 {
		return nil, dataset.ErrNoMore
	}

	if jns.current != nil {
		recs, err := jns.current.Pull(ctx, batch)
		if err != nil {
			if err != dataset.ErrNoMore {
				return nil, err
			}

			jns.current = nil
			if len(jns.streams) == 0 {
				return nil, dataset.ErrNoMore
			}

			goto NextStream
		}

		return recs, nil
	}

NextStream:
	if len(jns.streams) == 0 {
		return nil, dataset.ErrNoMore
	}

	next := jns.streams[0]
	jns.streams = jns.streams[1:]
	jns.current = &next

	recs, err := jns.current.Pull(ctx, batch)
	if err != nil {
		if err != dataset.ErrNoMore {
			return nil, err
		}

		goto NextStream
	}

	return recs, nil
}
