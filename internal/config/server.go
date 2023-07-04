package config

type Server struct {
	Host            string `env:"SERVER_HOST"`
	Port            int    `env:"SERVER_PORT"`
	ShutdownTimeout int    `env:"SERVER_SHUTDOWN_TIMEOUT"`
}
