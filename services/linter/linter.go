package linter

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// NodeCount is an analyzer that counts the number of each kind of node in a file.
var NodeCount = &analysis.Analyzer{
	Name: "nodecount",
	Doc:  "count the number of each kind of node in a file",
	Run:  runNodeCount,
}

// runNodeCount is the run function of the NodeCount analyzer.
func runNodeCount(pass *analysis.Pass) (interface{}, error) {
	// Create a map to store the node counts.
	counts := make(map[string]int)

	// Iterate over the files in the analysis unit.
	for _, file := range pass.Files {
		// Inspect the AST of the file.
		ast.Inspect(file, func(n ast.Node) bool {
			// If the node is nil, return false to stop the traversal.
			if n == nil {
				return false
			}

			// Get the name of the node type.
			name := fmt.Sprintf("%T", n)

			// Increment the count for the node type.
			counts[name]++

			// Return true to continue the traversal.
			return true
		})
	}

	// Report the node counts as a result.
	return counts, nil
}
