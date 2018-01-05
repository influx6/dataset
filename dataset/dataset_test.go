package dataset_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/influx6/faux/tests"
	"github.com/influx6/geckodataset/dataset"
)

func TestDataset(t *testing.T) {
	pusher := &mockaPush{}

	var set dataset.Dataset
	set.Pull = mockaPull{}
	set.Proc = mockaProc{}
	set.Pushers = append(set.Pushers, pusher)

	tests.Header("Should be able to process incoming records from source for storing")
	{
		pusher.Fn = func(recs ...map[string]interface{}) error {
			if len(recs) != 3 {
				return errors.New("expected to have being requested to pushed 3 records")
			}

			for index, rec := range recs {
				score, ok := rec["score"].(int)
				if !ok {
					return errors.New("expected to have 'score' field in processed record")
				}

				switch index {
				case 0:
					if score != 71 {
						tests.Failed("Should have received processed data with score of %d but got %d", 71, score)
					}
					tests.Passed("Should have received processed data with score of %d", 71)
				case 1:
					if score != 24 {
						tests.Failed("Should have received processed data with score of %d but got %d", 24, score)
					}
					tests.Passed("Should have received processed data with score of %d", 24)
				case 2:
					if score != 23 {
						tests.Failed("Should have received processed data with score of %d but got %d", 23, score)
					}
					tests.Passed("Should have received processed data with score of %d", 23)
				}
			}
			return nil
		}

		if err := set.Do(context.Background(), 100, 100); err != nil {
			tests.FailedWithError(err, "Should have successfully processed records")
		}
		tests.Passed("Should have successfully processed records")
	}

	tests.Header("Should be able to pull 3 records but push 1 at a time")
	{
		pusher.Fn = func(recs ...map[string]interface{}) error {
			if len(recs) != 1 {
				return fmt.Errorf("expected to have being requested to pushed %d records but got %d", 1, len(recs))
			}
			return nil
		}

		if err := set.Do(context.Background(), 3, 1); err != nil {
			tests.FailedWithError(err, "Should have successfully processed records")
		}
		tests.Passed("Should have successfully processed records")
	}

	tests.Header("Should be able to pull 2 records but push 2 at a time")
	{

		pusher.Fn = func(recs ...map[string]interface{}) error {
			if len(recs) != 2 {
				return fmt.Errorf("expected to have being requested to pushed %d records but got %d", 2, len(recs))
			}
			return nil
		}

		if err := set.Do(context.Background(), 2, 2); err != nil {
			tests.FailedWithError(err, "Should have successfully processed records")
		}
		tests.Passed("Should have successfully processed records")
	}

	tests.Header("Should be able to push all records with lower pull size")
	{

		pusher.Fn = func(recs ...map[string]interface{}) error {
			if len(recs) != 2 {
				return fmt.Errorf("expected to have being requested to pushed %d records but got %d", 2, len(recs))
			}
			return nil
		}

		if err := set.Do(context.Background(), 2, 3); err != nil {
			tests.FailedWithError(err, "Should have successfully processed records")
		}
		tests.Passed("Should have successfully processed records")
	}
}

type mockaPush struct {
	Fn func(...map[string]interface{}) error
}

func (m mockaPush) Push(ctx context.Context, recs ...map[string]interface{}) error {
	return m.Fn(recs...)
}

type mockaPull struct{}

func (m mockaPull) Pull(ctx context.Context, batch int) ([]map[string]interface{}, error) {
	recs := []map[string]interface{}{
		{
			"user":   "Rhis Whilly",
			"scores": []int{1, 32, 4, 5, 6, 23},
		},
		{
			"user":   "Josh Gambler",
			"scores": []int{1, 23},
		},
		{
			"user":   "Felix Decker",
			"scores": 23,
		},
	}

	if batch >= len(recs) {
		return recs, nil
	}

	return recs[:batch], nil
}

type mockaProc struct{}

func (m mockaProc) Transform(ctx context.Context, recs ...map[string]interface{}) ([]map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0, len(recs))
	for _, rec := range recs {
		res = append(res, map[string]interface{}{
			"user":  rec["name"],
			"score": getScore(rec["scores"]),
		})
	}
	return res, nil
}

func getScore(item interface{}) int {
	switch mo := item.(type) {
	case int:
		return mo
	case []int:
		var c int
		for _, elem := range mo {
			c += elem
		}
		return c
	default:
		return 0
	}
}
