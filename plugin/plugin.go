// Copyright (c) 2026 Michael_Go. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

// Package plugin provides the golangci-lint v2 module plugin entry point
// for the loglinter analyzer.
package plugin

import (
	"golang.org/x/tools/go/analysis"

	"github.com/dementievme/golangci-lint-plugin/internal/analyzer"
)

// plugin implements the golangci-lint AnalyzerPlugin interface.
type plugin struct{}

// GetAnalyzers returns the list of analyzers provided by this plugin.
func (plugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{analyzer.New()}
}

// AnalyzerPlugin is the exported symbol that golangci-lint discovers
// when loading this module as a custom linter plugin.
var AnalyzerPlugin plugin
