package utils

import (
	"fmt"
	"github.com/MikeMwita/go-strict/models"
	"github.com/MikeMwita/go-strict/services/complexity"
)

func PrintSummary(fileCount, funcCount, maxComplexity, totalComplexity, funcCountSum, complexLineCount int) {
	fmt.Printf("%d files\n%d functions\n%d highest complexity\n", fileCount, funcCount, maxComplexity)

	if funcCount == 0 {
		fmt.Println("No functions were found to calculate the average complexity per function")
	} else {
		avgComplexity := float64(totalComplexity) / float64(funcCountSum)
		fmt.Printf("%.2f overall average complexity per function\n", avgComplexity)
	}
	fmt.Printf("%d complex lines\n\n", complexLineCount)
}

func PrintDetails(results []*models.LintResult, format string, detailsFormat bool) {
	for _, result := range results {
		switch format {
		case "json":
			// JSON format printing logic
		case "line-number", "complexity":
			fmt.Printf("%s:%d:1 - %s has complexity: %s\n", result.File, result.Line, result.Function, result.Message)
		default:
			fmt.Printf("%s:%d:1 - %s has complexity: %s\n", result.File, result.Line, result.Function, result.Message)
		}

		if detailsFormat {
			details := complexity.GetDetail(result)
			fmt.Println("```go")
			fmt.Println(details)
			fmt.Println("```")
		}
	}
}
