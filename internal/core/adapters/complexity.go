package adapters

import (
	"go/ast"
	"go/token"

	"github.com/MikeMwita/go-strict/models"
)

type ComplexityCalculator interface {
	Calculate(fset *token.FileSet, file *ast.File, funcDecl *ast.FuncDecl) models.FunctionResult
}
