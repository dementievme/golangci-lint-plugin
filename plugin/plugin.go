package plugin

import (
	"golang.org/x/tools/go/analysis"

	"github.com/dementievme/golangci-lint-plugin/internal/analyzer"
	"github.com/dementievme/golangci-lint-plugin/internal/config"
)

type plugin struct{}

func (plugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{analyzer.New(config.Load())}
}

var AnalyzerPlugin plugin
