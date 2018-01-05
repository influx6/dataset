package jsonfiles_test

import (
	"context"
	"testing"

	"github.com/influx6/dataset/dataset"
	"github.com/influx6/dataset/dataset/pullers/jsonfiles"
	"github.com/influx6/faux/tests"
)

func TestJSONStreamWithBadJSONFile(t *testing.T) {
	jsx, err := jsonfiles.NewJSONStream("./fixtures/bad/bad.json")
	if err != nil {
		tests.FailedWithError(err, "Should have successfully loaded stream")
	}
	tests.Passed("Should have successfully loaded stream")

	_, err = jsx.Pull(context.Background(), 3)
	if err == nil {
		tests.Failed("Should have failed to load file content due to wrong format.")
	}
	tests.PassedWithError(err, "Should have failed to load file content due to wrong format.")

}

func TestJSONStreamWithGoodJSONFile(t *testing.T) {
	jsx, err := jsonfiles.NewJSONStream("./fixtures/sentos/rack.json")
	if err != nil {
		tests.FailedWithError(err, "Should have successfully loaded stream")
	}
	tests.Passed("Should have successfully loaded stream")

	docs, err := jsx.Pull(context.Background(), 3)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully pulled new records from source")
	}
	tests.Passed("Should have successfully pulled new records from source")

	if len(docs) != 3 {
		tests.Failed("Should have successfully retrieved 3 messages but got %d", len(docs))
	}
	tests.Passed("Should have successfully retrieved 3 messages")
}

func TestJSONStreamsWithFlatWalk(t *testing.T) {
	jssm, err := jsonfiles.New("./fixtures/sentos", false)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully loaded json files")
	}
	tests.Passed("Should have successfully loaded json files")

	if jssm.Total() != 3 {
		tests.Failed("Should have loaded 3 files %d", jssm.Total())
	}
	tests.Passed("Should have loaded 3 files")
}

func TestJSONStreamsWithDeepWalk(t *testing.T) {
	jssm, err := jsonfiles.New("./fixtures/sentos", true)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully loaded json files")
	}
	tests.Passed("Should have successfully loaded json files")

	if jssm.Total() != 4 {
		tests.Failed("Should have loaded 4 files %d", jssm.Total())
	}
	tests.Passed("Should have loaded 4 files")
}

func TestJSONStreamsLoad(t *testing.T) {
	tests.Header("When attempting to load files from root directory")
	{
		jssm, err := jsonfiles.New("./fixtures/sentos", false)
		if err != nil {
			tests.FailedWithError(err, "Should have successfully loaded json files")
		}
		tests.Passed("Should have successfully loaded json files")

		for {
			_, err := jssm.Pull(context.Background(), 2)
			if err != nil {
				// if we are at the end of the readlist, and streams has no more files, then we succeeded.
				if err == dataset.ErrNoMore && jssm.Total() == 0 {
					tests.Passed("Should have successfully read all data from files ")
					break
				}

				if err == dataset.ErrNoMore && jssm.Total() != 0 {
					tests.Failed("Should have successfully read all data from files ")
					break
				}

				tests.FailedWithError(err, "Should have successfully read next file content")
			}
		}
	}

	tests.Header("When attempting to load files from deep directory")
	{
		jssm, err := jsonfiles.New("./fixtures/sentos", true)
		if err != nil {
			tests.FailedWithError(err, "Should have successfully loaded json files")
		}
		tests.Passed("Should have successfully loaded json files")

		for {
			_, err := jssm.Pull(context.Background(), 2)
			if err != nil {
				// if we are at the end of the readlist, and streams has no more files, then we succeeded.
				if err == dataset.ErrNoMore && jssm.Total() == 0 {
					tests.Passed("Should have successfully read all data from files ")
					break
				}

				if err == dataset.ErrNoMore && jssm.Total() != 0 {
					tests.Failed("Should have successfully read all data from files ")
					break
				}

				tests.FailedWithError(err, "Should have successfully read next file content")
			}
		}
	}
}

func TestStreamRuns(t *testing.T) {
	for i := 0; i < 3; i++ {
		tests.Header("Should be able to get %d records from stream", i)
		{
			runStreams(t, i, func(recs []map[string]interface{}) {
				if len(recs) > i {
					tests.Failed("Should have received only a total of %d records but got %d", i, len(recs))
				}
			})
			tests.Passed("Should have received only a total of %d records", i)
		}
	}
}

func runStreams(t *testing.T, batch int, after func([]map[string]interface{})) {
	jssm, err := jsonfiles.New("./fixtures/sentos", false)
	if err != nil {
		tests.FailedWithError(err, "Should have successfully loaded json files")
	}
	tests.Passed("Should have successfully loaded json files")

	for {
		recs, err := jssm.Pull(context.Background(), batch)
		if err != nil {
			if err == dataset.ErrNoMore {
				break
			}

			tests.FailedWithError(err, "Should have successfully read next file content")
		}

		after(recs)
	}
}
