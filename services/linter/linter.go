package linter

import (
	"errors"
	"fmt"
	"github.com/MikeMwita/go-strict/datamodels"
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
	config     *datamodels.LintConfig
	complexity *complexity.ComplexityService
	fileCount  int
	funcCount  int
}

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
			// print the error
			fmt.Println("Error getting file info:", err)
			return nil, err
		}

		// check if the file is a directory
		if info.IsDir() {
			// get all the Go files in the directory
			goFiles, err := filepath.Glob(filepath.Join(file, "*.go"))
			if err != nil {
				fmt.Println("Error getting Go files in directory:", err)
				return nil, err
			}

			// print the files or directories that are passed as arguments
			fmt.Println("Linting files or directories:", files)

			fmt.Println("File info:", info)

			fmt.Println("Go files in the directory:", goFiles)

			// lint each Go file in the directory
			for _, goFile := range goFiles {
				// parse the Go file
				f, err := parser.ParseFile(fset, goFile, nil, parser.ParseComments)
				if err != nil {
					fmt.Println("Error parsing Go file:", err)

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

				// linting the Go file
				fileResults, err := ls.lintFile(fset, f)
				if err != nil {
					fmt.Println("Error linting file:", err)
					return nil, err
				}

				fmt.Println("File name:", goFile)

				// append the file results to the slice
				results = append(results, fileResults...)

				// update the file count
				fileCount++
			}
		} else {
			// parse the file
			f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
			if err != nil {
				fmt.Println("Error parsing Go file:", err)

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
				fmt.Println("Error linting file:", err)
				return nil, err
			}

			fmt.Println("File name:", file)

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

	// check if the results slice is empty
	if len(results) == 0 {
		// print a message indicating that no results were found
		fmt.Println("No results were found")
	} else {
		// print the results slice
		fmt.Println("Results:", results)
	}

	// return the linting results
	return results, nil
}

func (ls *LinterService) lintFile(fset *token.FileSet, f *ast.File) ([]*datamodels.LintResult, error) {
	// create a slice to store the file results
	var fileResults []*datamodels.LintResult

	// get the file name
	fileName := fset.File(f.Pos()).Name()

	// check if the file name is valid and not empty
	if fileName == "" {
		return nil, errors.New("empty file name")
	}

	// check if the file exists and is not a directory
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return nil, errors.New("file is a directory")
	}

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

	// increment the file count
	ls.fileCount++

	// return the file results only if they are not empty
	if len(fileResults) > 0 {
		return fileResults, nil
	}

	return nil, nil
}

// Inside lintFunction function
func (ls *LinterService) lintFunction(fset *token.FileSet, funcDecl *ast.FuncDecl) (*datamodels.LintResult, error) {
	if funcDecl.Body == nil {
		return nil, nil
	}

	complexity, err := ls.complexity.Calculate(fset, funcDecl.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Function: %s, Complexity: %d\n", funcDecl.Name.Name, complexity) // Added log statement

	if complexity > ls.config.Threshold {
		result := &datamodels.LintResult{
			Line:     fset.Position(funcDecl.Pos()).Line,
			Severity: "warning",
		}
		var details []string
		details = append(details, fmt.Sprintf("function has a cognitive complexity of %d which is higher than the threshold of %d", complexity, ls.config.Threshold))

		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
				line := fset.Position(node.Pos()).Line
				score := ls.complexity.Complexity(node)
				detail := fmt.Sprintf("+ %d (found at line: %d)", score, line)
				details = append(details, detail)
			case *ast.CaseClause:
				line := fset.Position(node.Pos()).Line
				score := ls.complexity.Complexity(node)
				detail := fmt.Sprintf("+ %d (found 'case' at line: %d)", score, line)
				details = append(details, detail)
			default:
				return true
			}
			return true
		})

		result.Message = fmt.Sprintf("%s\n%s", result.Message, strings.Join(details, "\n"))
		return result, nil
	}

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

	// writes the function declarations to the temporary file
	for _, function := range functions {
		if _, err := tmpFile.WriteString(function + "\n\n"); err != nil {
			return nil, err
		}
	}

	return tmpFile, nil
}

func NewLinterService(config *datamodels.LintConfig, complexity *complexity.ComplexityService) *LinterService {
	return &LinterService{
		config:     config,
		complexity: complexity,
	}
}
