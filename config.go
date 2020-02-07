package main

import (
	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
)

type Config struct {
	UrlDB   string
	UrlMem  string
	MemTTL  int32
	Workers int
}

const configFileName = "config.yml"

func initConfig() (*Config, error) {
	cfg := &Config{}

	if err := configor.New(&configor.Config{ENVPrefix: ""}).Load(cfg, configFileName); err != nil {
		return nil, errors.Wrap(err, "can't init configuration")
	}

	return cfg, nil
}
