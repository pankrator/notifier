package config

type SMSConfig struct {
	SID    string `env:"SMS_SID"`
	Secret string `env:"SMS_SECRET"`
	Sender string `env:"SMS_SENDER"`
}
