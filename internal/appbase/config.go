package appbase

import (
	"gr-scanner/pkg/config"
)

type Config struct {
	GitHubToken string `env:"GITHUB_TOKEN" env-required:"true"`
	LogLevel    string `env:"LOG_LEVEL" env-default:"debug"`
}

func LoadConfig() (*Config, error) {
	c := new(Config)

	err := config.LoadConfig(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
