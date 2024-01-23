package linter

import (
	"errors"
	"fmt"
	"github.com/MikeMwita/go-strict/models"
	"github.com/MikeMwita/go-strict/services/complexity"
	"github.com/MikeMwita/go-strict/utils"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Linter interface {
	LintFiles(files []string) ([]*models.LintResult, error)
	LintFunctions(functions []string) ([]*models.LintResult, error)
	lintFile(fset *token.FileSet, f *ast.File) ([]*models.LintResult, error)
	lintFunction(fset *token.FileSet, funcDecl *ast.FuncDecl) (*models.LintResult, error)
}

type LinterService struct {
	config     *models.LintConfig
	complexity *complexity.ComplexityService
	fileCount  int
	funcCount  int
}

func (ls *LinterService) LintFiles(files []string) ([]*models.LintResult, error) {
	var results []*models.LintResult
	fset := token.NewFileSet()

	for _, file := range files {
		err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error walking file tree: %v", err)
				return err
			}

			if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
				f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
				if err != nil {
					log.Printf("Error parsing Go file %s: %v", path, err)
					result := &models.LintResult{
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
			log.Printf("Error walking file tree: %v", err)
			return nil, err
		}
	}

	return results, nil
}

func (ls *LinterService) lintFile(fset *token.FileSet, f *ast.File) ([]*models.LintResult, error) {
	var fileResults []*models.LintResult

	fileName := fset.File(f.Pos()).Name()

	if fileName == "" {
		return nil, errors.New("empty file name")
	}

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return nil, errors.New("file is a directory")
	}

	for _, decl := range f.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			funcName := funcDecl.Name.Name

			// lint the function
			funcResult, err := ls.lintFunction(fset, funcDecl)
			if err != nil {
				return nil, err
			}

			if funcResult != nil {
				funcResult.File = fileName
				funcResult.Function = funcName

				fileResults = append(fileResults, funcResult)
			}
		}
	}

	ls.fileCount++

	if len(fileResults) > 0 {
		return fileResults, nil
	}

	return nil, nil
}

// lints the given function declaration

func (ls *LinterService) lintFunction(fset *token.FileSet, funcDecl *ast.FuncDecl) (*models.LintResult, error) {
	if funcDecl.Body == nil {
		return nil, nil
	}

	complexityService := complexity.NewComplexityService()
	complexity, err := complexityService.Calculate(fset, funcDecl.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Function: %s, Complexity: %d\n", funcDecl.Name.Name, complexity)

	if complexity > ls.config.Threshold {
		result := &models.LintResult{
			Line:     fset.Position(funcDecl.Pos()).Line,
			Severity: "warning",
		}
		var details []string
		details = append(details, fmt.Sprintf("function has a cognitive complexity of %d which is higher than the threshold of %d", complexity, ls.config.Threshold))

		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
				line := fset.Position(node.Pos()).Line
				score := complexityService.Complexity(node)
				detail := fmt.Sprintf("+ %d (found at line: %d)", score, line)
				details = append(details, detail)
			case *ast.CaseClause:
				line := fset.Position(node.Pos()).Line
				score := complexityService.Complexity(node)
				detail := fmt.Sprintf("+ %d (found 'case' at line: %d)", score, line)
				details = append(details, detail)
			default:
				return true
			}
			return true
		})

		result.Message = fmt.Sprintf("%s (Complexity details:\n%s)", result.Message, strings.Join(details, "\n"))

		results := []*models.LintResult{result}

		utils.PrintComplexity(results, "line-number", true)
		fmt.Printf("Complexity details for function %s in file %s at line %d: %d\n", funcDecl.Name.Name, fset.Position(funcDecl.Pos()).Filename, fset.Position(funcDecl.Pos()).Line, complexity)

		return result, nil
	}
	return nil, nil
}

func (ls *LinterService) LintFunctions(functions []string) ([]*models.LintResult, error) {
	var results []*models.LintResult

	fset := token.NewFileSet()

	tmpFile, err := createTempFile(functions)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	f, err := parser.ParseFile(fset, tmpFile.Name(), nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, decl := range f.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			funcResult, err := ls.lintFunction(fset, funcDecl)
			if err != nil {
				fmt.Printf("Error linting function %s: %v\n", funcDecl.Name.Name, err)
				continue
			}

			if funcResult != nil {
				results = append(results, funcResult)
			}
		} else {
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

	for _, function := range functions {
		if _, err := tmpFile.WriteString(function + "\n\n"); err != nil {
			return nil, err
		}
	}

	return tmpFile, nil
}

func NewLinterService(config *models.LintConfig, complexity *complexity.ComplexityService) *LinterService {
	return &LinterService{
		config:     config,
		complexity: complexity,
	}
}
