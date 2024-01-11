package complexity

import (
	"fmt"
	"go/ast"
	"go/token"
)

type ComplexityService struct {
	nesting int
}

// calculates complexity for any statement(general)

func (cs *ComplexityService) Complexity(stmt ast.Node) int {
	// Complexity calculation for any statement
	return 1 + cs.nesting
}

// calculates the complexity of a function body

func (cs *ComplexityService) Calculate(fset *token.FileSet, node ast.Node) (int, error) {

	// Initializing complexity score
	complexity := 1 // Starting with 1 as the base complexity

	// Traversing the function body to calculate complexity
	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.IfStmt:
			// increment the complexity for the if statement
			complexity += cs.Complexity(n)

			// increment the nesting level for the if statement
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
			// increment the complexity for the loop statement
			complexity += cs.Complexity(n)

			// increment the nesting level for the loop statement
			cs.nesting++

			// traverse the loop statement body, key, value, and expression
			// NB: we are using  n directly without type assertion
			complexity, _ = cs.Calculate(fset, n.Body)
			if n.Key != nil {
				complexity, _ = cs.Calculate(fset, n.Key)
			}
			if n.Value != nil {
				complexity, _ = cs.Calculate(fset, n.Value)
			}
			complexity, _ = cs.Calculate(fset, n.X)

			// decrement the nesting level after the loop statement
			cs.nesting--
		case *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
			// increment the complexity for the switch statement
			complexity += cs.Complexity(n)

			// increment the nesting level for the switch statement
			cs.nesting++

			// traverse the switch statement body, init, and tag
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

			// decrement the nesting level after the switch statement
			cs.nesting--
		case *ast.CaseClause:
			// increment the complexity for the case clause
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

//  prints the complexity of a function

func PrintComplexity(cs *ComplexityService, fset *token.FileSet, fn *ast.FuncDecl) {
	// call the Calculate method of the ComplexityService type
	complexity, err := cs.Calculate(fset, fn.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	// print the function name and complexity
	fmt.Printf("%s has complexity: %d\n", fn.Name, complexity)
}

func NewComplexityService() *ComplexityService {
	return &ComplexityService{
		nesting: 0,
	}
}
