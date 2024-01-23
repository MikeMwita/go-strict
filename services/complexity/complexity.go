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

func (cs *ComplexityService) Complexity(stmt ast.Node) int {
	return 1 + cs.nesting
}

// function body complexity calculation

func (cs *ComplexityService) Calculate(fset *token.FileSet, node ast.Node) (int, error) {

	complexity := 1

	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.IfStmt:
			// increment the complexity for the if statement
			complexity += cs.Complexity(n)

			cs.nesting++

			// traverse the if statement body, init, and else branches
			complexity, _ = cs.Calculate(fset, n.Body)
			if n.Init != nil {
				complexity, _ = cs.Calculate(fset, n.Init)
			}
			if n.Else != nil {
				complexity, _ = cs.Calculate(fset, n.Else)
			}

			// decrement the nesting level after the if statement
			cs.nesting--
		case *ast.ForStmt:
			// increment the complexity for the loop statement
			complexity += cs.Complexity(n)

			// increment the nesting level for the loop statement
			cs.nesting++

			// traverse the loop statement body, init, post, and condition

			complexity, _ = cs.Calculate(fset, n.Body)
			if n.Init != nil {
				complexity, _ = cs.Calculate(fset, n.Init)
			}
			if n.Post != nil {
				complexity, _ = cs.Calculate(fset, n.Post)
			}
			if n.Cond != nil {
				complexity, _ = cs.Calculate(fset, n.Cond)
			}

			// decrement the nesting level after the loop statement
			cs.nesting--
		case *ast.RangeStmt:
			complexity += cs.Complexity(n)

			cs.nesting++

			complexity, _ = cs.Calculate(fset, n.Body)
			if n.Key != nil {
				complexity, _ = cs.Calculate(fset, n.Key)
			}
			if n.Value != nil {
				complexity, _ = cs.Calculate(fset, n.Value)
			}
			complexity, _ = cs.Calculate(fset, n.X)

			cs.nesting--
		case *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
			complexity += cs.Complexity(n)

			cs.nesting++

			switch n := n.(type) {
			case *ast.SwitchStmt:
				complexity, _ = cs.Calculate(fset, n.Body)
				if n.Init != nil {
					complexity, _ = cs.Calculate(fset, n.Init)
				}
				if n.Tag != nil {
					complexity, _ = cs.Calculate(fset, n.Tag)
				}
			case *ast.TypeSwitchStmt:
				complexity, _ = cs.Calculate(fset, n.Body)
				if n.Init != nil {
					complexity, _ = cs.Calculate(fset, n.Init)
				}
				complexity, _ = cs.Calculate(fset, n.Assign)
			case *ast.SelectStmt:
				complexity, _ = cs.Calculate(fset, n.Body)
			}

			cs.nesting--
		case *ast.CaseClause:
			complexity += cs.Complexity(n)

			// create a BlockStmt to represent the case clause body
			block := &ast.BlockStmt{
				Lbrace: n.Colon,
				Rbrace: n.Colon,
				List:   n.Body,
			}

			// traverse the case clause body
			complexity, _ = cs.Calculate(fset, block)
		}
		return true
	})

	return complexity, nil
}

func GetDetail(result *models.LintResult) string {
	var details string

	details += fmt.Sprintf(" (Complexity details:\n%s)", result.Message)

	return details
}

func NewComplexityService() *ComplexityService {
	return &ComplexityService{
		nesting: 0,
	}
}
