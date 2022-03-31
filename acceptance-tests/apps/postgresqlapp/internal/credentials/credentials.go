package credentials

import (
	"fmt"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	URI              string `mapstructure:"uri"`
	ClientCACert     string `mapstructure:"client_ca_cert"`
	ClientPrivateKey string `mapstructure:"client_private_key"`
	Hostname         string `mapstructure:"hostname"`
}

func Read() (*Config, error) {
	app, err := cfenv.Current()
	if err != nil {
		return nil, fmt.Errorf("error reading app env: %w", err)
	}
	svs, err := app.Services.WithTag("postgresql")
	if err != nil {
		return nil, fmt.Errorf("error reading PostgreSQL service details")
	}
	c := &Config{}

	if err := mapstructure.Decode(svs[0].Credentials, c); err != nil {
		return nil, fmt.Errorf("failed to decode credentials: %w", err)
	}

	if c.URI == "" {
		return nil, fmt.Errorf("parsed credentials are not valid")
	}

	fmt.Printf("config:: %v\n", c)

	return c, nil
}
