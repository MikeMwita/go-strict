package adapters

import (
	"go/ast"
	"go/token"
)

type Complexity interface {
	Calculate(fset *token.FileSet, node ast.Node) (int, error)
	//GetDetail(result *models.LintResult) string
}

type ComplexityCalculator interface {
	Calculate(fset *token.FileSet, body *ast.BlockStmt) (int, error)
	If(stmt *ast.IfStmt) int
	Loop(stmt ast.Node) int
	Switch(stmt ast.Node) int
	Case(stmt *ast.CaseClause) int
}
