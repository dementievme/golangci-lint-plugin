package validator

import (
	"github.com/dementievme/golangci-lint-plugin/internal/config"
)

type Rule func(msg string) error

type Validator struct {
	rules []Rule
}

var defaultSensitiveKeywords = []string{
	"password", "pwd", "pass", "token", "token_secret", "pwd", "bearer",
}

func NewValidator(config *config.Config) *Validator {
	keywords := append(defaultSensitiveKeywords, config.ExtraSensitiveKeywords...)

	disabled := make(map[string]bool)
	for _, r := range config.DisableRules {
		disabled[r] = true
	}

	v := &Validator{}

	if !disabled["lowercase"] {
		v.rules = append(v.rules, Lowercase())
	}
	if !disabled["special_chars"] {
		v.rules = append(v.rules, SpecialChars())
	}
	if !disabled["english"] {
		v.rules = append(v.rules, English())
	}
	if !disabled["sensitive_data"] {
		v.rules = append(v.rules, SensitiveData(keywords))
	}

	return v
}

func (v *Validator) Validate(msg string) []error {
	var errs []error
	for _, rule := range v.rules {
		if err := rule(msg); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
