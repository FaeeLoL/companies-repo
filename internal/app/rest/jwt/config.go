package jwt

import (
	"github.com/spf13/viper"
)

type Config struct {
	SecretKey      string
	TrustedIssuers []string
}

func NewConfig() *Config {
	return &Config{}
}

func LoadJWTConfig() *Config {
	viper.SetDefault("server.jwt.secret_key", "secret")
	viper.SetDefault("server.jwt.trusted_issuers", []string{"trusted.issuer"})

	return &Config{
		SecretKey:      viper.GetString("server.jwt.secret_key"),
		TrustedIssuers: viper.GetStringSlice("server.jwt.trusted_issuers"),
	}
}
