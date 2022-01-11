package credentials

import (
	"fmt"
	"log"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/mitchellh/mapstructure"
)

type BigQueryCredentials struct {
	Credentials string `mapstructure:"credentials"`
	ProjectID   string `mapstructure:"projectId"`
	DatasetID   string `mapstructure:"dataset_id"`
}

func Read() (BigQueryCredentials, error) {
	app, err := cfenv.Current()
	if err != nil {
		return BigQueryCredentials{}, fmt.Errorf("error reading app env: %w", err)
	}
	svs, err := app.Services.WithTag("bigquery")
	if err != nil {
		return BigQueryCredentials{}, fmt.Errorf("error reading BigQuery service details")
	}

	var r BigQueryCredentials

	if err := mapstructure.Decode(svs[0].Credentials, &r); err != nil {
		return BigQueryCredentials{}, fmt.Errorf("failed to decode credentials: %w", err)
	}

	log.Println("creds: %s", r)
	if r.Credentials == "" || r.ProjectID == "" || r.DatasetID == "" {
		return BigQueryCredentials{}, fmt.Errorf("parsed credentials are not valid")
	}

	return r, nil
}
