package linter

import (
	"fmt"
	"github.com/MikeMwita/go-strict/models"
	"github.com/MikeMwita/go-strict/services/complexity"
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
}

type LinterService struct {
	config     *models.LintConfig
	complexity *complexity.ComplexityService
	fileCount  int
	funcCount  int
}

func NewLinterService(config *models.LintConfig, complexity *complexity.ComplexityService) *LinterService {
	return &LinterService{
		config:     config,
		complexity: complexity,
	}
}

func createTempFile(functions []string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile("", "tempfunctions*.go")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}

	for _, function := range functions {
		if _, err := tmpFile.WriteString(function + "\n\n"); err != nil {
			return nil, fmt.Errorf("failed to write to temporary file: %w", err)
		}
	}

	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temporary file: %w", err)
	}

	return tmpFile, nil
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
				fileResults, err := ls.lintGoFile(fset, path)
				if err != nil {
					log.Printf("Error linting Go file %s: %v", path, err)
					return err
				}
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
				log.Printf("Error linting function %s: %v", funcDecl.Name.Name, err)
				continue
			}

			if funcResult != nil {
				results = append(results, funcResult)
			}
		}
	}

	return results, nil
}

func (ls *LinterService) lintGoFile(fset *token.FileSet, filePath string) ([]*models.LintResult, error) {
	f, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Printf("Error parsing Go file %s: %v", filePath, err)
		return nil, err
	}

	return ls.lintFile(fset, f)
}

func (ls *LinterService) lintFile(fset *token.FileSet, f *ast.File) ([]*models.LintResult, error) {
	var fileResults []*models.LintResult
	fileName := fset.File(f.Pos()).Name()

	for _, decl := range f.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			funcResult, err := ls.lintFunction(fset, funcDecl)
			if err != nil {
				return nil, err
			}

			if funcResult != nil {
				funcResult.File = fileName
				fileResults = append(fileResults, funcResult)
			}
		}
	}

	ls.fileCount++
	return fileResults, nil
}

func (ls *LinterService) lintFunction(fset *token.FileSet, funcDecl *ast.FuncDecl) (*models.LintResult, error) {
	if funcDecl.Body == nil {
		return nil, nil
	}

	complexityScore, err := ls.complexity.Calculate(fset, funcDecl.Body)
	if err != nil {
		return nil, err
	}

	if complexityScore > ls.config.Threshold {
		result := &models.LintResult{
			File:     fset.Position(funcDecl.Pos()).Filename,
			Line:     fset.Position(funcDecl.Pos()).Line,
			Function: funcDecl.Name.Name,
			Message:  fmt.Sprintf("function has a cognitive complexity of %d which is higher than the threshold of %d", complexityScore, ls.config.Threshold),
			Severity: "warning",
		}

		details := ls.generateComplexityDetails(fset, funcDecl.Body)
		result.Message = fmt.Sprintf("%s (Complexity details:\n%s)", result.Message, strings.Join(details, "\n"))

		return result, nil
	}
	return nil, nil
}

func (ls *LinterService) generateComplexityDetails(fset *token.FileSet, body *ast.BlockStmt) []string {
	var details []string
	ast.Inspect(body, func(n ast.Node) bool {
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
		}
		return true
	})
	return details
}
