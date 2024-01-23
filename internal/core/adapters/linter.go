package adapters

import (
	"github.com/MikeMwita/go-strict/models"
	"go/ast"
	"go/token"
)

type Linter interface {
	LintFiles(files []string) ([]*models.LintResult, error)
	LintFunctions(functions []string) ([]*models.LintResult, error)
	lintFile(fset *token.FileSet, f *ast.File) ([]*models.LintResult, error)
	lintFunction(fset *token.FileSet, funcDecl *ast.FuncDecl) (*models.LintResult, error)
}
