package config

import "github.com/joeshaw/envdecode"

type Config struct {
	Server        Server        `env:""`
	EmailerConfig EmailerConfig `env:""`
	SlackerConfig SlackerConfig `env:""`
	SMSConfig     SMSConfig     `env:""`
	DB            DB            `env:""`
	Processor     Processor     `env:""`
}

func NewConfig() (*Config, error) {
	c := &Config{}

	if err := envdecode.Decode(c); err != nil {
		return nil, err
	}

	return c, nil
}
