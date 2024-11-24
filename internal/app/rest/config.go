package rest

import (
	"time"

	"github.com/spf13/viper"

	"github.com/faeelol/companies-store/internal/app/rest/jwt"
)

type Config struct {
	Addr              string
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	IdleTimeout       time.Duration
	JWT               *jwt.Config
}

func NewConfig() *Config {
	return &Config{
		JWT: jwt.NewConfig(),
	}
}

func LoadServerConfig() *Config {
	viper.SetDefault("server.addr", ":8080")
	viper.SetDefault("server.write_timeout", "15s")
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.read_header_timeout", "5s")
	viper.SetDefault("server.idle_timeout", "60s")

	return &Config{
		Addr:              viper.GetString("server.addr"),
		WriteTimeout:      viper.GetDuration("server.write_timeout"),
		ReadTimeout:       viper.GetDuration("server.read_timeout"),
		ReadHeaderTimeout: viper.GetDuration("server.read_header_timeout"),
		IdleTimeout:       viper.GetDuration("server.idle_timeout"),
		JWT:               jwt.LoadJWTConfig(),
	}
}
