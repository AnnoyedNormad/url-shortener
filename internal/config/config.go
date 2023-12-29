package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env            string `yaml:"env" env-default:"development"`
	StoragePath    string `yaml:"storage_path" env-required:"true"`
	HTTPServer     `yaml:"http_server"`
	User           `yaml:"user"`
	AliasGenerator `yaml:"alias_generator"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type User struct {
	Login    string `yaml:"login" env-required:"true"`
	Password string `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type AliasGenerator struct {
	Length   int    `yaml:"length" env-default:"6"`
	Alphabet string `yaml:"alphabet" env-default:"qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890"`
	GenLimit int    `yaml:"gen_limit" env-default:"10"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", configPath)
	}

	return &cfg
}
