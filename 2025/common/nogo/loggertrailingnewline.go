package loggertrailingnewline

import (
	"fmt"
	"go/ast"
	"go/token"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "loggertrailingnewline",
	Doc:  "Detects logger calls that include trailing newlines",
	Run:  run,
}

var loggingFunctions = []string{"Infof", "Errorf", "Fatalf"}

func checkFile(pass *analysis.Pass, f *ast.File) {
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		lhs, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		if lhs.Name != "logger" {
			return true
		}

		ident := sel.Sel
		if ident == nil || !slices.Contains(loggingFunctions, ident.Name) {
			return true
		}

		if len(call.Args) < 1 {
			pass.Reportf(call.Pos(), "no args")
			return true
		}

		msg, ok := call.Args[0].(*ast.BasicLit)
		if !ok {
			return true
		}

		if msg.Kind != token.STRING {
			return true
		}

		if strings.HasSuffix(msg.Value, "\\n\"") {
			pass.Reportf(call.Pos(),
				fmt.Sprintf("logger.%v format with trailing newline", ident.Name))
		}

		return true
	})
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		checkFile(pass, f)
	}
	return nil, nil
}
