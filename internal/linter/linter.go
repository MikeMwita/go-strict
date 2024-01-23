package linter

import (
	"errors"
	"fmt"
	"github.com/MikeMwita/go-strict/datamodels"
	"github.com/MikeMwita/go-strict/services/complexity"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Linter interface {
	LintFiles(files []string) ([]*datamodels.LintResult, error)
	LintFunctions(functions []string) ([]*datamodels.LintResult, error)
	lintFile(fset *token.FileSet, f *ast.File) ([]*datamodels.LintResult, error)
	lintFunction(fset *token.FileSet, funcDecl *ast.FuncDecl) (*datamodels.LintResult, error)
}

type LinterService struct {
	config     *datamodels.LintConfig
	complexity *complexity.ComplexityService
	fileCount  int
	funcCount  int
}

func (ls *LinterService) LintFiles(files []string) ([]*datamodels.LintResult, error) {
	var results []*datamodels.LintResult

	fset := token.NewFileSet()

	for _, file := range files {
		err := fs.WalkDir(os.DirFS(file), ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Printf("Error walking file tree: %v", err)
				return err
			}

			if strings.HasSuffix(d.Name(), ".go") {
				f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
				if err != nil {
					log.Printf("Error parsing Go file %s: %v", path, err)

					result := &datamodels.LintResult{
						File:     path,
						Message:  err.Error(),
						Severity: "error",
					}

					results = append(results, result)

					return nil
				}

				fileResults, err := ls.lintFile(fset, f)
				if err != nil {
					log.Printf("Error linting file %s: %v", path, err)
					return err
				}

				log.Printf("Linting file: %s", path)

				results = append(results, fileResults...)
			}

			return nil
		})

		if err != nil {
			// log the error
			log.Printf("Error walking file tree: %v", err)
			return nil, err
		}
	}

	// return the results and nil error
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

	fmt.Printf("Function: %s, Complexity: %d\n", funcDecl.Name.Name, complexity)

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

		// Include complexity details in the result message
		result.Message = fmt.Sprintf("%s (Complexity details:\n%s)", result.Message, strings.Join(details, "\n"))

		cs := complexity.NewComplexityService()
		complexityDetails, _ := cs.Calculate(fset, funcDecl.Body)
		complexityFn := funcDecl.Name.Name
		complexityFile := fset.Position(funcDecl.Pos()).Filename
		complexityFnLine := fset.Position(funcDecl.Pos()).Line
		complexityComplexity := complexityDetails
		complexity.PrintComplexity(cs, fset, funcDecl)
		fmt.Printf("Complexity details for function %s in file %s at line %d: %d\n", complexityFn, complexityFile, complexityFnLine, complexityComplexity)

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
