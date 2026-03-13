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

var defaultLoggers = map[string]map[string]bool{
	"log/slog": {
		"Debug": true, "Info": true, "Warn": true, "Error": true,
		"DebugContext": true, "InfoContext": true, "WarnContext": true, "ErrorContext": true,
	},
	"go.uber.org/zap": {
		"Debug": true, "Info": true, "Warn": true, "Error": true,
		"Fatal": true, "Panic": true,
		"Debugf": true, "Infof": true, "Warnf": true, "Errorf": true,
		"Fatalf": true, "Panicf": true,
		"Debugw": true, "Infow": true, "Warnw": true, "Errorw": true,
		"Fatalw": true, "Panicw": true,
	},
	"log": {
		"Print": true, "Printf": true, "Println": true,
		"Fatal": true, "Fatalf": true, "Fatalln": true,
		"Panic": true, "Panicf": true, "Panicln": true,
	},
}

var defaultSensitiveKeywords = []string{
	"password", "pwd", "pass", "token", "token_secret", "pwd", "bearer",
}

func Load() *Config {
	cfg := &Config{}

	if path, err := fetchConfigPath(); err == nil {
		_ = cleanenv.ReadConfig(path, cfg)
	}

	cfg.applyDefaults()
	return cfg
}

func (c *Config) applyDefaults() {
	if len(c.Loggers) == 0 {
		c.Loggers = defaultLoggers
	}
	if len(c.ExtraSensitiveKeywords) == 0 {
		c.ExtraSensitiveKeywords = defaultSensitiveKeywords
	}
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
