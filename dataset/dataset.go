package dataset

import (
	"context"
)

// Proc defines an interface which embodies a a processor of
// records.
type Proc interface {
	Transform(context.Context, ...map[string]interface{}) ([]map[string]interface{}, error)
}

// DataPull defines an interface which exposes a pull method to
// collect specific amount of records from underline store.
type DataPull interface {
	Pull(context.Context, int) ([]map[string]interface{}, error)
}

// DataPush embodies a interface exposing a push method to
// store incoming records of map instances into underline store.
type DataPush interface {
	Push(context.Context, ...map[string]interface{}) error
}

// DataPushers implements the DataPush for a slice of DataPush items
// where each is called with provided map records.
type DataPushers []DataPush

// Push runs all Pusher within slice type and returns when a pusher meets
// an error or when all pushers have successfully pushed provided map records.
func (dp DataPushers) Push(ctx context.Context, recs ...map[string]interface{}) error {
	for _, pusher := range dp {
		if err := pusher.Push(ctx, recs...); err != nil {
			return err
		}
	}
	return nil
}

// Dataset implements a custom data processor which takes implementations
// of the DataPull and DataPush(optional) interfaces, where the provided Procs
// instance processes data received from the Pull and stored into the Push
// implementation.
type Dataset struct {
	Proc    Proc
	Pull    DataPull
	Pushers DataPushers
}

// Do immediately runs the conversion process to transform data received from
// the puller into the pushers list.
// Do is to be used recursively, where every call processes the next batch taking
// from the the puller and processed, if an error occured, then that error will be
// returned.
func (ds Dataset) Do(ctx context.Context, pullBatch int, pushBatch int) error {
	recs, err := ds.Pull.Pull(ctx, pullBatch)
	if err != nil {
		return err
	}

	procRecs, err := ds.Proc.Transform(ctx, recs...)
	if err != nil {
		return err
	}

	if pushBatch >= len(procRecs) || pullBatch < pushBatch {
		return ds.Pushers.Push(ctx, procRecs...)
	}

	totalBatch := len(recs) / pushBatch
	for i := 1; i <= totalBatch; i++ {
		if pushBatch >= len(recs) {
			if err = ds.Pushers.Push(ctx, recs...); err != nil {
				return err
			}
			return nil
		}

		next := recs[:pushBatch]
		recs = recs[pushBatch:]
		if err = ds.Pushers.Push(ctx, next...); err != nil {
			return err
		}
	}

	return nil
}
