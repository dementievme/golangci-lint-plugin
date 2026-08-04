// Copyright (c) 2026 Michael_Go. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

// Package config provides configuration loading for the loglinter analyzer.
// It supports reading settings from a YAML file and merging user-defined
// sensitive keywords with built-in defaults.
package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config holds the loglinter configuration.
type Config struct {
	// Loggers maps package import paths to a set of method names to check.
	Loggers map[string]map[string]bool `yaml:"loggers"`
	// ExtraSensitiveKeywords is a list of additional keywords that indicate
	// sensitive data. These are merged with the built-in defaults.
	ExtraSensitiveKeywords []string `yaml:"extra_sensitive_keywords"`
	// DisableRules is a list of rule names to disable (e.g. "lowercase", "english").
	DisableRules []string `yaml:"disable_rules"`
}

// defaultLoggers defines the built-in set of supported loggers and their methods.
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

// defaultSensitiveKeywords is the built-in list of keywords that may indicate
// sensitive data exposure in log messages.
var defaultSensitiveKeywords = []string{
	"password", "pwd", "pass", "token", "token_secret", "bearer",
}

// Load reads loglinter configuration from the file at path.
// If path is empty, it falls back to the CONFIG_PATH environment variable.
// If the config file cannot be read, a warning is logged and defaults are used.
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

// applyDefaults fills in default values for unset fields and merges
// user-provided extra sensitive keywords with the built-in list.
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
