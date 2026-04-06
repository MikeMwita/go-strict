// Package analyzer exposes a go/analysis.Analyzer so go-strict can be used as
// a golangci-lint plugin or run directly with the `go vet -vettool` mechanism.
//
// Usage with golangci-lint custom plugin:
//
//	linters-settings:
//	  custom:
//	    go-strict:
//	      path: ./go-strict.so
//	      description: Cognitive complexity linter
//	      settings:
//	        threshold: 10
package analyzer

import (
	"flag"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/MikeMwita/go-strict/config"
	"github.com/MikeMwita/go-strict/internal/complexity"
)

var Analyzer = &analysis.Analyzer{
	Name:     "gocognitive",
	Doc:      "Reports functions whose cognitive complexity exceeds the configured threshold.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

var threshold int

func init() {
	Analyzer.Flags = *flag.NewFlagSet("gocognitive", flag.ContinueOnError)
	Analyzer.Flags.IntVar(&threshold, "threshold", 10,
		"cognitive complexity threshold above which functions are reported")
}

func run(pass *analysis.Pass) (interface{}, error) {
	cfg := config.DefaultConfig()
	if threshold > 0 {
		cfg.Threshold = threshold
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			return
		}

		// Find the file this function belongs to.
		var file *ast.File
		for _, f := range pass.Files {
			if f.Pos() <= funcDecl.Pos() && funcDecl.End() <= f.End() {
				file = f
				break
			}
		}

		calc := complexity.NewCalculator(pass.Fset, file, cfg)
		result := calc.Calculate(funcDecl)

		if result.Exceeded {
			pass.Reportf(
				funcDecl.Pos(),
				"function %s has cognitive complexity %d (threshold %d)",
				result.Function,
				result.Complexity,
				cfg.Threshold,
			)
		}
	})

	return nil, nil
}
