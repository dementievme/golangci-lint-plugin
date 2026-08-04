package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
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
	"password", "pwd", "pass", "token", "token_secret", "bearer",
}

func Load(path string) *Config {
	cfg := &Config{}

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	if path != "" {
		if err := cleanenv.ReadConfig(path, cfg); err != nil {
			log.Printf("loglinter: warning: failed to read config %s: %v", path, err)
		}
	}

	cfg.applyDefaults()
	return cfg
}

func (cfg *Config) applyDefaults() {
	if len(cfg.Loggers) == 0 {
		cfg.Loggers = defaultLoggers
	}

	seen := make(map[string]bool)
	var merged []string
	for _, kw := range defaultSensitiveKeywords {
		if !seen[kw] {
			seen[kw] = true
			merged = append(merged, kw)
		}
	}
	for _, kw := range cfg.ExtraSensitiveKeywords {
		if !seen[kw] {
			seen[kw] = true
			merged = append(merged, kw)
		}
	}
	cfg.ExtraSensitiveKeywords = merged
}
