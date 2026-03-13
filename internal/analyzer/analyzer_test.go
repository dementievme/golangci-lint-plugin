package analyzer_test

import (
	"testing"

	"github.com/dementievme/golangci-lint-plugin/internal/analyzer"
	"github.com/dementievme/golangci-lint-plugin/internal/config"
	"golang.org/x/tools/go/analysis/analysistest"
)

func testConfig() *config.Config {
	return &config.Config{
		Loggers: map[string]map[string]bool{
			"log/slog": {
				"Info": true, "Debug": true, "Warn": true, "Error": true,
			},
		},
		ExtraSensitiveKeywords: []string{"password", "token"},
	}
}

func TestAnalyzer(t *testing.T) {
	a := analyzer.New(testConfig())
	analysistest.Run(t, analysistest.TestData(), a, "examples")
}
