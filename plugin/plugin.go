package plugin

import (
	"golang.org/x/tools/go/analysis"

	"github.com/dementievme/golangci-lint-plugin/internal/analyzer"
)

type plugin struct{}

func (plugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{analyzer.New()}
}

var AnalyzerPlugin plugin
