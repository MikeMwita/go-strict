package complexity

import (
	"fmt"
	"github.com/MikeMwita/go-strict/models"
	"go/ast"
	"go/token"
)

type ComplexityService struct {
	nesting int
}

func NewComplexityService() *ComplexityService {
	return &ComplexityService{
		nesting: 0,
	}
}

func (cs *ComplexityService) Calculate(fset *token.FileSet, node ast.Node) (int, error) {
	complexity := 1
	ast.Inspect(node, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
			complexity += cs.Complexity(n)
		case *ast.CaseClause:
			complexity += cs.Complexity(n)
		}
		return true
	})

	return complexity, nil
}

func (cs *ComplexityService) Complexity(stmt ast.Node) int {
	return 1 + cs.nesting
}

func GetDetail(result *models.LintResult) string {
	var details string
	details += fmt.Sprintf(" (Complexity details:\n%s)", result.Message)
	return details
}
