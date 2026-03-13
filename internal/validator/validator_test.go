package validator_test

import (
	"testing"

	"github.com/dementievme/golangci-lint-plugin/internal/config"
	"github.com/dementievme/golangci-lint-plugin/internal/validator"
)

func cfg(disable ...string) *config.Config {
	return &config.Config{
		DisableRules:           disable,
		ExtraSensitiveKeywords: []string{"password", "token"},
	}
}

func TestLowercase(t *testing.T) {
	v := validator.New(cfg("special_chars", "english", "sensitive_data"))

	cases := []struct {
		msg  string
		fail bool
	}{
		{"starting server", false},
		{"Starting server", true},
		{"ALLCAPS", true},
		{"", false},
	}

	for _, tc := range cases {
		errs := v.Validate(tc.msg)
		if tc.fail && len(errs) == 0 {
			t.Errorf("expected error for %q", tc.msg)
		}
		if !tc.fail && len(errs) != 0 {
			t.Errorf("unexpected error for %q: %v", tc.msg, errs)
		}
	}
}

func TestEnglish(t *testing.T) {
	v := validator.New(cfg("lowercase", "special_chars", "sensitive_data"))

	cases := []struct {
		msg  string
		fail bool
	}{
		{"starting server", false},
		{"запуск сервера", true},
		{"server started", false},
		{"ошибка подключения", true},
	}

	for _, tc := range cases {
		errs := v.Validate(tc.msg)
		if tc.fail && len(errs) == 0 {
			t.Errorf("expected error for %q", tc.msg)
		}
		if !tc.fail && len(errs) != 0 {
			t.Errorf("unexpected error for %q: %v", tc.msg, errs)
		}
	}
}

func TestSpecialChars(t *testing.T) {
	v := validator.New(cfg("lowercase", "english", "sensitive_data"))

	cases := []struct {
		msg  string
		fail bool
	}{
		{"server started", false},
		{"server started!", true},
		{"connection failed!!!", true},
		{"something went wrong", false},
	}

	for _, tc := range cases {
		errs := v.Validate(tc.msg)
		if tc.fail && len(errs) == 0 {
			t.Errorf("expected error for %q", tc.msg)
		}
		if !tc.fail && len(errs) != 0 {
			t.Errorf("unexpected error for %q: %v", tc.msg, errs)
		}
	}
}

func TestSensitiveData(t *testing.T) {
	v := validator.New(cfg("lowercase", "english", "special_chars"))

	cases := []struct {
		msg  string
		fail bool
	}{
		{"user authenticated", false},
		{"user password exposed", true},
		{"token validated", true},
		{"api request completed", false},
	}

	for _, tc := range cases {
		errs := v.Validate(tc.msg)
		if tc.fail && len(errs) == 0 {
			t.Errorf("expected error for %q", tc.msg)
		}
		if !tc.fail && len(errs) != 0 {
			t.Errorf("unexpected error for %q: %v", tc.msg, errs)
		}
	}
}

func TestDisableRules(t *testing.T) {
	v := validator.New(cfg("lowercase", "english", "special_chars", "sensitive_data"))
	errs := v.Validate("Starting server! запуск token")
	if len(errs) != 0 {
		t.Errorf("expected no errors with all rules disabled, got: %v", errs)
	}
}

func TestMultipleErrors(t *testing.T) {
	v := validator.New(cfg())
	errs := v.Validate("Starting server!")
	if len(errs) < 2 {
		t.Errorf("expected at least 2 errors, got %d", len(errs))
	}
}
