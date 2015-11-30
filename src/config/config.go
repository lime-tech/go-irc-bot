package config

import (
	"github.com/BurntSushi/toml"
	"go-irc-bot/src/client"
)

type Config struct {
	Clients map[string]*client.Config
	Http    *HttpApi
}
type HttpApi struct {
	Addr        string
	MaxPostSize int64
}

func FromFile(path string) (*Config, error) {
	config := new(Config)
	config.Http = &HttpApi{
		Addr:        "127.0.0.1:8080",
		MaxPostSize: 1048576,
	}

	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}

	return config, nil
}
