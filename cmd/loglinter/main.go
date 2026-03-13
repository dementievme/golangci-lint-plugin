package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/dementievme/golangci-lint-plugin/internal/analyzer"
	"github.com/dementievme/golangci-lint-plugin/internal/config"
)

func main() {
	config := config.Load()
	analyzer := analyzer.New(config)

	singlechecker.Main(analyzer)
}
