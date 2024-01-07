package linter

import (
	"fmt"
	"github.com/MikeMwita/go-strict/domain/datamodels"
	"github.com/MikeMwita/go-strict/services/complexity"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// Linter is an interface that defines the linting methods
type Linter interface {
	LintFiles(files []string) ([]*datamodels.LintResult, error)         // lints the given files or directories
	LintFunctions(functions []string) ([]*datamodels.LintResult, error) // lints the given functions
}

// LinterService is a service that implements the Linter interface
type LinterService struct {
	config     *datamodels.LintConfig // the linting configuration
	complexity complexity.Complexity  // the complexity calculator
}

// NewLinterService creates a new LinterService
func NewLinterService(config *datamodels.LintConfig, complexity complexity.Complexity) *LinterService {
	return &LinterService{
		config:     config,
		complexity: complexity,
	}
}

// LintFiles lints the given files or directories
func (ls *LinterService) LintFiles(files []string) ([]*datamodels.LintResult, error) {
	// create a slice to store the linting results
	var results []*datamodels.LintResult

	// create a file set to parse the files
	fset := token.NewFileSet()

	// iterate over the files or directories
	for _, file := range files {
		// get the file info
		info, err := os.Stat(file)
		if err != nil {
			return nil, err
		}

		// check if the file is a directory
		if info.IsDir() {
			// get all the Go files in the directory
			goFiles, err := filepath.Glob(filepath.Join(file, "*.go"))
			if err != nil {
				return nil, err
			}

			// lint each Go file in the directory
			for _, goFile := range goFiles {
				// parse the Go file
				f, err := parser.ParseFile(fset, goFile, nil, parser.ParseComments)
				if err != nil {
					return nil, err
				}

				// lint the Go file
				fileResults, err := ls.lintFile(fset, f)
				if err != nil {
					return nil, err
				}

				// append the file results to the slice
				results = append(results, fileResults...)
			}
		} else {
			// parse the file
			f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
			if err != nil {
				return nil, err
			}

			// lint the file
			fileResults, err := ls.lintFile(fset, f)
			if err != nil {
				return nil, err
			}

			// append the file results to the slice
			results = append(results, fileResults...)
		}
	}

	// return the linting results
	return results, nil
}

// lintFile lints a single file
func (ls *LinterService) lintFile(fset *token.FileSet, f *ast.File) ([]*datamodels.LintResult, error) {
	// create a slice to store the file results
	var fileResults []*datamodels.LintResult

	// get the file name
	fileName := fset.File(f.Pos()).Name()

	// iterate over the declarations in the file
	for _, decl := range f.Decls {
		// check if the declaration is a function declaration
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// get the function name
			funcName := funcDecl.Name.Name

			// lint the function
			funcResult, err := ls.lintFunction(fset, funcDecl)
			if err != nil {
				return nil, err
			}

			// check if the function result is not nil
			if funcResult != nil {
				// set the file name and the function name
				funcResult.File = fileName
				funcResult.Function = funcName

				// append the function result to the file results
				fileResults = append(fileResults, funcResult)
			}
		}
	}

	// return the file results
	return fileResults, nil
}

// lintFunction lints a single function
func (ls *LinterService) lintFunction(fset *token.FileSet, funcDecl *ast.FuncDecl) (*datamodels.LintResult, error) {
	// check if the function has a body
	if funcDecl.Body == nil {
		return nil, nil
	}

	// calculate the complexity of the function
	complexity, err := ls.complexity.Calculate(fset, funcDecl.Body)
	if err != nil {
		return nil, err
	}

	// check if the complexity exceeds the threshold
	if complexity > ls.config.Threshold {
		// create a lint result
		result := &datamodels.LintResult{
			Line:     fset.Position(funcDecl.Pos()).Line,
			Message:  fmt.Sprintf("function has a cognitive complexity of %d which is higher than the threshold of %d", complexity, ls.config.Threshold),
			Severity: "warning",
		}

		// return the lint result
		return result, nil
	}

	// return nil if the complexity is within the threshold
	return nil, nil
}

func (ls *LinterService) LintFunctions(functions []string) (interface{}, interface{}) {

	return nil, nil
}

// NodeCount is an analyzer that counts the number of each kind of node in a file.
//var NodeCount = &analysis.Analyzer{
//	Name: "nodecount",
//	Doc:  "count the number of each kind of node in a file",
//	Run:  runNodeCount,
//}
//
//// runNodeCount is the run function of the NodeCount analyzer.
//func runNodeCount(pass *analysis.Pass) (interface{}, error) {
//	// Create a map to store the node counts.
//	counts := make(map[string]int)
//
//	// Iterate over the files in the analysis unit.
//	for _, file := range pass.Files {
//		// Inspect the AST of the file.
//		ast.Inspect(file, func(n ast.Node) bool {
//			// If the node is nil, return false to stop the traversal.
//			if n == nil {
//				return false
//			}
//
//			// Get the name of the node type.
//			name := fmt.Sprintf("%T", n)
//
//			// Increment the count for the node type.
//			counts[name]++
//
//			// Return true to continue the traversal.
//			return true
//		})
//	}

//	// Report the node counts as a result.
//	return counts, nil
//}
