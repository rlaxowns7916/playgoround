package config

import (
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Kafka    KafkaConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Database string
	URL      string
	Username string
	Password string
}

type KafkaConfig struct {
	BootstrapServers []string
}

func New() (*Config, error) {
	k := koanf.New(".")
	if err := k.Load(file.Provider("./configs/config.yaml"), yaml.Parser()); err != nil {
		return nil, err
	}
	if err := k.Load(env.Provider("APP_", ".", func(s string) string {
		key := strings.TrimPrefix(s, "APP_")
		key = strings.ToLower(key)
		key = strings.ReplaceAll(key, "_", ".")
		key = strings.ReplaceAll(key, "bootstrap.servers", "bootstrap-servers")
		return key
	}), nil); err != nil {
		return nil, err
	}

	cfg := &Config{
		Server: ServerConfig{
			Host: k.String("server.host"),
			Port: k.Int("server.port"),
		},
		Database: DatabaseConfig{
			Database: k.String("database.database"),
			URL:      k.String("database.url"),
			Username: k.String("database.username"),
			Password: k.String("database.password"),
		},
		Kafka: KafkaConfig{
			BootstrapServers: k.Strings("kafka.bootstrap-servers"),
		},
	}

	return cfg, nil
}
