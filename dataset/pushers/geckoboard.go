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

	set.UniqueBy = conf.UniqueBy
	if err := client.Create(context.Background(), conf.Dataset, set); err != nil {
		return GeckoboardPusher{}, err
	}

	return GeckoboardPusher{
		Config: conf,
		Client: client,
	}, nil
}

// Replace takes incoming map of records which will be the transformed data received
// from the a Proc.
func (gh GeckoboardPusher) Update(ctx context.Context, recs ...map[string]interface{}) error {
	return gh.Client.ReplaceData(ctx, gh.Config.Dataset, geckoclient.Dataset{
		Data:     recs,
		DeleteBy: gh.Config.DeteletBy,
	})
}

// Push takes incoming map of records which will be the transformed data received
// from the a Proc.
func (gh GeckoboardPusher) Add(ctx context.Context, recs ...map[string]interface{}) error {
	return gh.Client.PushData(ctx, gh.Config.Dataset, geckoclient.Dataset{
		Data: recs,
	})
}

// Send uses the operation flag from the config to send giving records to the Geckoboard's dataset API.
func (gh GeckoboardPusher) Push(ctx context.Context, recs ...map[string]interface{}) error {
	switch strings.ToLower(gh.Config.Op) {
	case "push":
		return gh.Add(ctx, recs...)
	case "update":
		return gh.Update(ctx, recs...)
	default:
		return errors.New("unknown operation type")
	}
}

func transformFields(fields []config.FieldType) (geckoclient.NewDataset, error) {
	var set geckoclient.NewDataset
	set.Fields = map[string]geckoclient.DataType{}

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
