package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Loggers                map[string]map[string]bool `yaml:"loggers"`
	ExtraSensitiveKeywords []string                   `yaml:"extra_sensitive_keywords"`
	DisableRules           []string                   `yaml:"disable_rules"`
}

func MustLoad() *Config {
	config := &Config{}

	configPath, err := fetchConfigPath()
	if err != nil {
		panic(err)
	}

	if err := cleanenv.ReadConfig(configPath, config); err != nil {
		panic(err)
	}

	return config
}

func fetchConfigPath() (string, error) {
	_ = godotenv.Load()

	var configPath string
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	if configPath != "" {
		return configPath, nil
	}

	configPath, found := os.LookupEnv("CONFIG_PATH")
	if !found {
		return "", ErrPathNotSpecified
	}
	if configPath == "" {
		return "", ErrConfigPathIsEmpty
	}

	return configPath, nil
}
