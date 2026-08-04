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

// New creates a new loglinter analyzer with lazy config loading.
// Config path is read from the -config flag or CONFIG_PATH env variable.
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

// NewWithConfig creates a new loglinter analyzer with the given config.
// Intended for use in tests.
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

// messageArg extracts the message argument from a log call.
// For *Context methods (e.g. slog.InfoContext), the message is at index 1 (after ctx).
// For all other methods, the message is at index 0.
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
