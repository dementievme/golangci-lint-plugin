package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/dementievme/golangci-lint-plugin/internal/analyzer"
)

func main() {
	singlechecker.Main(analyzer.New())
}
