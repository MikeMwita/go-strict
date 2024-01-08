package linter

import (
	"fmt"
	"github.com/MikeMwita/go-strict/domain/datamodels"
	"github.com/MikeMwita/go-strict/services/complexity"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

	// create some variables to keep track of the statistics
	var fileCount, funcCount, totalComplexity, maxComplexity, complexLineCount int

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
					// create a lint result with the file error
					result := &datamodels.LintResult{
						File:     goFile,
						Message:  err.Error(),
						Severity: "error",
					}

					// append the result to the slice
					results = append(results, result)

					// continue with the next file
					continue
				}

				// lint the Go file
				fileResults, err := ls.lintFile(fset, f)
				if err != nil {
					return nil, err
				}

				// append the file results to the slice
				results = append(results, fileResults...)

				// update the file count
				fileCount++
			}
		} else {
			// parse the file
			f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
			if err != nil {
				// create a lint result with the file error
				result := &datamodels.LintResult{
					File:     file,
					Message:  err.Error(),
					Severity: "error",
				}

				// append the result to the slice
				results = append(results, result)

				// continue with the next file
				continue
			}

			// lint the file
			fileResults, err := ls.lintFile(fset, f)
			if err != nil {
				return nil, err
			}

			// append the file results to the slice
			results = append(results, fileResults...)

			// update the file count
			fileCount++
		}
	}

	// calculate the average complexity
	avgComplexity := float64(totalComplexity) / float64(funcCount)

	// format the statistics as a string
	stats := fmt.Sprintf("%d = files\n%d = functions\n%d = highest complexity\n%.2f = overall average complexity per function\n%d = complex lines\n", fileCount, funcCount, maxComplexity, avgComplexity, complexLineCount)

	// create a lint result with the statistics
	result := &datamodels.LintResult{
		Message: stats,
	}

	// append the result to the slice
	results = append(results, result)

	// return the linting results
	return results, nil
}

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
			Severity: "warning",
		}

		// create a slice to store the complexity details
		var details []string

		// format the complexity score for the function
		details = append(details, fmt.Sprintf("function has a cognitive complexity of %d which is higher than the threshold of %d", complexity, ls.config.Threshold))

		// traverse the function body and get the complexity details for each statement
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			// check the type and value of the node
			switch node := n.(type) {
			case *ast.IfStmt:
				// get the line number and the complexity score for the if statement
				line := fset.Position(node.Pos()).Line
				score := ls.complexity.If(node)

				// format the complexity detail for the if statement
				detail := fmt.Sprintf("+ %d (found 'if' at line: %d)", score, line)

				// append the detail to the slice
				details = append(details, detail)
			case *ast.ForStmt, *ast.RangeStmt:
				// get the line number and the complexity score for the loop statement
				line := fset.Position(node.Pos()).Line
				score := ls.complexity.Loop(node)

				// format the complexity detail for the loop statement
				detail := fmt.Sprintf("+ %d (found 'loop' at line: %d)", score, line)

				// append the detail to the slice
				details = append(details, detail)
			case *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
				// get the line number and the complexity score for the switch statement
				line := fset.Position(node.Pos()).Line
				score := ls.complexity.Switch(node)

				// format the complexity detail for the switch statement
				detail := fmt.Sprintf("+ %d (found 'switch' at line: %d)", score, line)

				// append the detail to the slice
				details = append(details, detail)
			case *ast.CaseClause:
				// get the line number and the complexity score for the case clause
				line := fset.Position(node.Pos()).Line
				score := ls.complexity.Case(node)

				// format the complexity detail for the case clause
				detail := fmt.Sprintf("+ %d (found 'case' at line: %d)", score, line)

				// append the detail to the slice
				details = append(details, detail)
			// Add more cases for other statement types as needed

			default:
				return true
			}

			return true
		})

		// append the complexity details to the lint result
		result.Message = fmt.Sprintf("%s\n%s", result.Message, strings.Join(details, "\n"))

		// return the lint result
		return result, nil
	}

	// return nil if the complexity is within the threshold
	return nil, nil
}

// LintFunctions lints the given functions
func (ls *LinterService) LintFunctions(functions []string) ([]*datamodels.LintResult, error) {
	// create a slice to store the linting results
	var results []*datamodels.LintResult

	// create a file set to parse the functions
	fset := token.NewFileSet()

	// create a temporary file with the provided functions
	tmpFile, err := createTempFile(functions)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	// parse the temporary file
	f, err := parser.ParseFile(fset, tmpFile.Name(), nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// iterate over the declarations in the file
	for _, decl := range f.Decls {
		// check if the declaration is a function declaration
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// lint the function
			funcResult, err := ls.lintFunction(fset, funcDecl)
			if err != nil {
				// log the error and proceed to the next function
				fmt.Printf("Error linting function %s: %v\n", funcDecl.Name.Name, err)
				continue
			}

			// check if the function result is not nil
			if funcResult != nil {
				// append the function result to the slice
				results = append(results, funcResult)
			}
		} else {
			// log a warning for non-function declarations
			fmt.Printf("Warning: unexpected declaration type %T\n", decl)
		}
	}

	// return the linting results
	return results, nil
}

// createTempFile creates a temporary file with the provided function declarations
func createTempFile(functions []string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile("", "tempfunctions*.go")
	if err != nil {
		return nil, err
	}
	defer tmpFile.Close()

	// write the function declarations to the temporary file
	for _, function := range functions {
		if _, err := tmpFile.WriteString(function + "\n\n"); err != nil {
			return nil, err
		}
	}

	return tmpFile, nil
}
