package complexity

import (
	"go/ast"
	"go/token"
)

// Complexity is an interface that defines the complexity calculation method
type Complexity interface {
	Calculate(fset *token.FileSet, body *ast.BlockStmt) (int, error) // calculates the complexity of a function body
	If(stmt *ast.IfStmt) int                                         // calculates complexity for if statements
	Loop(stmt ast.Node) int                                          // calculates complexity for loops
	Switch(stmt ast.Node) int                                        // calculates complexity for switch statements
	Case(stmt *ast.CaseClause) int                                   // calculates complexity for case clauses
}

// ComplexityService is a service that implements the Complexity interface
type ComplexityService struct {
	// you can add any fields or dependencies here
	nesting int // the current nesting level
}

// NewComplexityService creates a new ComplexityService
func NewComplexityService() *ComplexityService {
	return &ComplexityService{
		// you can initialize any fields or dependencies here
		nesting: 0, // starting with zero nesting level
	}
}

// Calculate calculates the complexity of a function body
func (cs *ComplexityService) Calculate(fset *token.FileSet, body *ast.BlockStmt) (int, error) {
	// you can implement the complexity calculation logic here
	// you can use the fset to get the file and line information
	// you can use the body to traverse the abstract syntax tree of the function body

	// Initializing complexity score
	complexity := 1 // Starting with 1 as the base complexity

	// Traversing the function body to calculate complexity
	ast.Inspect(body, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.IfStmt:
			// increment the complexity for the if statement
			complexity += cs.If(node)

			// increment the nesting level for the if statement
			cs.nesting++

			// traverse the if statement body
			ast.Inspect(node.Body, cs.calculate)

			// decrement the nesting level after the if statement
			cs.nesting--

			// check if the if statement has an else branch
			if node.Else != nil {
				// increment the complexity for the else branch
				complexity++

				// traverse the else branch
				ast.Inspect(node.Else, cs.calculate)
			}
		case *ast.ForStmt, *ast.RangeStmt:
			// increment the complexity for the loop statement
			complexity += cs.Loop(node)

			// increment the nesting level for the loop statement
			cs.nesting++

			// traverse the loop statement body
			//ast.Inspect(node.Body, cs.calculate)

			// decrement the nesting level after the loop statement
			cs.nesting--
		case *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
			// increment the complexity for the switch statement
			complexity += cs.Switch(node)

			// increment the nesting level for the switch statement
			cs.nesting++

			// traverse the switch statement body
			//ast.Inspect(node.Body, cs.calculate)

			// decrement the nesting level after the switch statement
			cs.nesting--
		case *ast.CaseClause:
			// increment the complexity for the case clause
			complexity += cs.Case(node)
			// create a BlockStmt to represent the case clause body
			block := &ast.BlockStmt{
				Lbrace: node.Colon,
				Rbrace: node.Colon,
				List:   node.Body,
			}

			// traverse the case clause body
			ast.Inspect(block, cs.calculate)
		}
		return true
	})

	return complexity, nil
}

// calculate is a helper function to calculate the complexity of a node
func (cs *ComplexityService) calculate(n ast.Node) bool {
	// Convert n to ast.Stmt using a type assertion
	stmt, ok := n.(ast.Stmt)
	if !ok {
		// n is not an ast.Stmt, do not include it in the calculation
		return true
	}

	// call the Calculate method with a dummy file set and block statement
	_, err := cs.Calculate(&token.FileSet{}, &ast.BlockStmt{List: []ast.Stmt{stmt}})
	return err == nil
	//// call the Calculate method with a dummy file set and block statement
	//_, err := cs.Calculate(&token.FileSet{}, &ast.BlockStmt{List: []ast.Stmt{n}})
	//return err == nil
}

// If calculates complexity for if statements
func (cs *ComplexityService) If(stmt *ast.IfStmt) int {
	// Complexity calculation for an if statement
	return 1 + cs.nesting
}

// Loop calculates complexity for loops
func (cs *ComplexityService) Loop(stmt ast.Node) int {
	// Complexity calculation for a loop statement
	return 1 + cs.nesting
}

// Switch calculates complexity for switch statements
func (cs *ComplexityService) Switch(stmt ast.Node) int {
	// Complexity calculation for a switch statement
	return 1 + cs.nesting
}

// Case calculates complexity for case clauses
func (cs *ComplexityService) Case(stmt *ast.CaseClause) int {
	// Complexity calculation for a case clause
	return 1 + cs.nesting
}
