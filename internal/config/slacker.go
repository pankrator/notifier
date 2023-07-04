package config

type SlackerConfig struct {
	Webhook string `env:"SLACK_WEBHOOK"`
}
