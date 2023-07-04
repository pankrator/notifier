package config

import "time"

type Processor struct {
	FetchInterval time.Duration `env:"PROCESSOR_FETCH_INTERVAL"`
	BatchSize     int           `env:"PROCESSOR_BATCH_SIZE"`
}
