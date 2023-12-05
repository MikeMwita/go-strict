package a

import "fmt"

func main() {
	var x int // want "found *ast.DeclStmt"
	x = 42 // want "found *ast.AssignStmt"
	fmt.Println(x) // want "found *ast.ExprStmt"
}