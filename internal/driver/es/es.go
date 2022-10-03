package es

import (
	"context"

	el "github.com/olivere/elastic/v7"
	"gitlab.com/stoqu/stoqu-be/internal/config"
	"gitlab.com/stoqu/stoqu-be/pkg/util/gracefull"
)

var (
	Client  *el.Client
	err     error
	stopper gracefull.ProcessStopper
)

func Init() (gracefull.ProcessStopper, error) {
	Client, err = el.NewClient(
		el.SetURL(config.Config.Driver.Elasticsearch.Url),
		el.SetSniff(false),
		// el.SetBasicAuth(config.Config.Driver.Elasticsearch.User, config.Config.Driver.Elasticsearch.Password),
	)
	if err != nil {
		return gracefull.NilStopper, err
	}
	stopper = func(ctx context.Context) (err error) {
		Client.Stop()
		return nil
	}
	return stopper, nil
}
