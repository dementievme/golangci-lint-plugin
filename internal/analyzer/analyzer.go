package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"

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

func New(config *config.Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     analyzerName,
		Doc:      analyzerDoc,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      run(validator.New(config), config.Loggers),
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

			msg, pos, ok := firstStringArg(call)
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

func firstStringArg(call *ast.CallExpr) (string, token.Pos, bool) {
	for _, arg := range call.Args {
		lit, ok := arg.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			continue
		}

		val, err := strconv.Unquote(lit.Value)
		if err != nil {
			continue
		}

		return val, lit.Pos(), true
	}

	return "", 0, false
}
