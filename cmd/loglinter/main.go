// Copyright (c) 2026 Michael_Go. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

// Command loglinter runs the loglinter analyzer as a standalone tool.
//
// Usage:
//
//	loglinter -config ./config/config.yml ./...
//	CONFIG_PATH=./config/config.yml loglinter ./...
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/dementievme/golangci-lint-plugin/internal/analyzer"
)

func main() {
	singlechecker.Main(analyzer.New())
}
