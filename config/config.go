package config

import (
	"github.com/BurntSushi/toml"
)

type HttpApi struct {
	Addr        string
	MaxPostSize int64
}

type Config struct {
	Server         string
	Nick           string
	Channels       []string
	ServerPassword string
	Http           *HttpApi
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
