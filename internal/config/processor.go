package config

import "time"

type Processor struct {
	FetchInterval time.Duration `env:"PROCESSOR_FETCH_INTERVAL"`
}
