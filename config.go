package main

import (
	"fmt"

	"github.com/knadh/koanf/v2"
)

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

type VoskConfig struct {
	Host string
	Port int
}

type AppConfig struct {
	Environment Environment
}

type Config struct {
	App  AppConfig
	Vosk VoskConfig
}

func getConfig(k *koanf.Koanf) *Config {
	c := Config{}

	c.App.Environment = Environment(fmt.Sprint(k.Get("app.environment")))
	c.Vosk.Host = fmt.Sprint(k.Get("vosk.host").(string))
	c.Vosk.Port = k.Get("vosk.port").(int)

	return &c
}
