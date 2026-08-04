// Copyright (c) 2026 Michael_Go. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

// Package analyzer implements the loglinter analysis pass.
// It inspects Go source files for log function calls and validates
// their message arguments against configurable style and security rules.
package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/dementievme/golangci-lint-plugin/internal/config"
	"github.com/dementievme/golangci-lint-plugin/internal/validator"
)

const (
	analyzerName = "loglinter"
	analyzerDoc  = "Checks log messages for style and security violations."
)

// New creates a loglinter analyzer with lazy config loading.
// The config path is read from the -config flag or the CONFIG_PATH
// environment variable at first run.
func New() *analysis.Analyzer {
	var configPath string
	a := &analysis.Analyzer{
		Name:     analyzerName,
		Doc:      analyzerDoc,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
	a.Flags.StringVar(&configPath, "config", "", "path to loglinter config file")
	a.Run = makeRun(&configPath)
	return a
}

// NewWithConfig creates a loglinter analyzer with the given config.
// This constructor is intended for use in tests where the config
// is known at compile time.
func NewWithConfig(cfg *config.Config) *analysis.Analyzer {
	v := validator.New(cfg)
	loggers := cfg.Loggers
	return &analysis.Analyzer{
		Name:     analyzerName,
		Doc:      analyzerDoc,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      run(v, loggers),
	}
}

// makeRun returns an analysis run function that lazily loads the config
// on its first invocation using sync.Once.
func makeRun(configPath *string) func(*analysis.Pass) (interface{}, error) {
	var once sync.Once
	var v *validator.Validator
	var loggers map[string]map[string]bool

	return func(pass *analysis.Pass) (interface{}, error) {
		once.Do(func() {
			cfg := config.Load(*configPath)
			v = validator.New(cfg)
			loggers = cfg.Loggers
		})

		return run(v, loggers)(pass)
	}
}

// run returns an analysis function that walks call expressions,
// filters log calls, extracts the message argument and validates it.
func run(v *validator.Validator, loggers map[string]map[string]bool) func(*analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		insp.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
			call := n.(*ast.CallExpr)

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || !isLogCall(pass, sel, loggers) {
				return
			}

			msg, pos, ok := messageArg(call, sel.Sel.Name)
			if !ok {
				return
			}

			for _, err := range v.Validate(msg) {
				pass.Reportf(pos, "%s", err.Error())
			}
		})

		return nil, nil
	}
}

// isLogCall reports whether sel refers to a known log function
// by resolving its type information and checking against the loggers map.
func isLogCall(pass *analysis.Pass, sel *ast.SelectorExpr, loggers map[string]map[string]bool) bool {
	obj, ok := pass.TypesInfo.Uses[sel.Sel]
	if !ok {
		return false
	}

	fn, ok := obj.(*types.Func)
	if !ok || fn.Pkg() == nil {
		return false
	}

	return loggers[fn.Pkg().Path()][sel.Sel.Name]
}

// messageArg extracts the message string literal from a log call.
// For *Context methods (e.g. slog.InfoContext), the message is at
// argument index 1 (after ctx). For all other methods, it is at index 0.
// Returns false if the argument is not a string literal.
func messageArg(call *ast.CallExpr, methodName string) (string, token.Pos, bool) {
	idx := 0
	if strings.HasSuffix(methodName, "Context") {
		idx = 1
	}
	if idx >= len(call.Args) {
		return "", 0, false
	}

	lit, ok := call.Args[idx].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return "", 0, false
	}

	val, err := strconv.Unquote(lit.Value)
	if err != nil {
		return "", 0, false
	}

	return val, lit.Pos(), true
}
