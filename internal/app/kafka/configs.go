package kafka

import (
	"github.com/spf13/viper"
)

type Config struct {
	Brokers []string
	Topic   string
}

func LoadKafkaConfig() *Config {
	return &Config{
		Brokers: viper.GetStringSlice("kafka.brokers"),
		Topic:   viper.GetString("kafka.topic"),
	}
}
