package pushers

import (
	"context"
	"errors"

	"strings"

	"github.com/influx6/geckoclient"
	"github.com/influx6/geckodataset/dataset/config"
)

// GeckoboardPusher implements the Pusher interface for sending data to the
// Geckoboard API for the user's account identified by the auth key.
type GeckoboardPusher struct {
	created bool
	Client  geckoclient.Client
	Config  config.DatasetConfig
	NewSet  *geckoclient.NewDataset
}

// NewGeckoboardPusher returns a new instance of GeckoboardPusher.
func NewGeckoboardPusher(apiKey string, conf config.DatasetConfig) (GeckoboardPusher, error) {
	client, err := geckoclient.New(apiKey)
	if err != nil {
		return GeckoboardPusher{}, err
	}

	// transform fields to dataset record.
	set, err := transformFields(conf.Fields)
	if err != nil {
		return GeckoboardPusher{}, err
	}

	return GeckoboardPusher{
		Config: conf,
		Client: client,
		NewSet: &set,
	}, nil
}

// FindOrCreate attempts to create pushers dataset if the pusher has an associated
// NewDataSet field value which indicates need to create set if not existing.
func (gh GeckoboardPusher) FindOrCreate(ctx context.Context) error {
	if gh.NewSet == nil {
		return nil
	}
	return gh.Client.Create(ctx, gh.Config.Dataset, *gh.NewSet)
}

// Push takes incoming map of records which will be the transformed data received
// from the a Proc.
func (gh GeckoboardPusher) Push(ctx context.Context, recs ...map[string]interface{}) error {
	return gh.Client.ReplaceData(ctx, gh.Config.Dataset, geckoclient.Dataset{
		Data: recs,
	})
}

func transformFields(fields []config.FieldType) (geckoclient.NewDataset, error) {
	var set geckoclient.NewDataset

	for _, desc := range fields {
		if desc.Name == "" {
			return set, errors.New("Name value is required for dataset field")
		}

		if desc.Type == "" {
			return set, errors.New("Type value is required for dataset field")
		}

		switch strings.ToLower(desc.Type) {
		case "date":
			set.Fields[desc.Name] = geckoclient.DateType{
				Name: desc.Name,
			}
		case "money":
			if desc.Currency == "" {
				return set, errors.New("Currency value is required for Money dataset field")
			}

			set.Fields[desc.Name] = geckoclient.MoneyType{
				Name:         desc.Name,
				CurrencyCode: desc.Currency,
				Optional:     desc.Optional,
			}
		case "string":
			set.Fields[desc.Name] = geckoclient.StringType{
				Name: desc.Name,
			}
		case "number":
			set.Fields[desc.Name] = geckoclient.NumberType{
				Name:     desc.Name,
				Optional: desc.Optional,
			}
		case "datetime":
			set.Fields[desc.Name] = geckoclient.DateTimeType{
				Name: desc.Name,
			}
		case "percentage":
			set.Fields[desc.Name] = geckoclient.PercentageType{
				Name:     desc.Name,
				Optional: desc.Optional,
			}
		}
	}
	return set, nil
}
