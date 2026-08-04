// Copyright (c) 2026 Michael_Go. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

// Package validator provides log message validation against
// a configurable set of style and security rules.
package validator

import (
	"github.com/dementievme/golangci-lint-plugin/internal/config"
)

// Rule is a function that validates a log message.
// It returns an error describing the violation, or nil if the message is valid.
type Rule func(msg string) error

// Validator applies a set of rules to log messages.
type Validator struct {
	rules []Rule
}

// New creates a Validator with rules determined by the given config.
// Rules listed in cfg.DisableRules are excluded.
func New(cfg *config.Config) *Validator {
	keywords := cfg.ExtraSensitiveKeywords

	disabled := make(map[string]bool)
	for _, r := range cfg.DisableRules {
		disabled[r] = true
	}

	v := &Validator{}

	if !disabled["lowercase"] {
		v.rules = append(v.rules, Lowercase())
	}

	if !disabled["english"] {
		v.rules = append(v.rules, English())
	}

	if !disabled["special_chars"] {
		v.rules = append(v.rules, SpecialChars())
	}

	if !disabled["sensitive_data"] {
		v.rules = append(v.rules, SensitiveData(keywords))
	}

	return v
}

// Validate runs all enabled rules against the given message
// and returns a slice of validation errors (may be empty).
func (v *Validator) Validate(msg string) []error {
	var errs []error
	for _, rule := range v.rules {
		if err := rule(msg); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
