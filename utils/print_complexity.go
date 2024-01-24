package utils

import (
	"fmt"
	"github.com/MikeMwita/go-strict/models"
	"github.com/MikeMwita/go-strict/services/complexity"
)

type ComplexityService struct {
	nesting int
}

func PrintComplexity(results []*models.LintResult, format string, detailsFormat bool) {
	var fileCount, funcCount, totalComplexity, maxComplexity, complexLineCount int

	for _, result := range results {
		if result.File != "" {
			fileCount++
		}

		if result.Function != "" {
			funcCount++

			var complexity int
			fmt.Sscanf(result.Message, "function has a cognitive complexity of %d", &complexity)

			totalComplexity += complexity
			if complexity > maxComplexity {
				maxComplexity = complexity
			}

			// our assumed threshold is 10,
			if complexity > 12 {
				complexLineCount++
			}
		}
	}

	fmt.Printf("%d = files\n%d = functions\n%d = highest complexity\n", fileCount, funcCount, maxComplexity)

	if funcCount == 0 {
		fmt.Println("No functions were found to calculate the average complexity per function")
	} else {
		avgComplexity := float64(totalComplexity) / float64(funcCount)
		fmt.Printf("%.2f = overall average complexity per function\n", avgComplexity)
	}

	fmt.Printf("%d = complex lines\n\n", complexLineCount)

	// Check if the format is "complexity"
	if format == "complexity" {
		// Change the format to "line-number"
		format = "line-number"
	}

	for _, result := range results {
		if result.File != "" && result.Function != "" {
			printResultDetails(result, format, detailsFormat)
		}
	}
}

func printResultDetails(result *models.LintResult, format string, detailsFormat bool) {
	if format == "line-number" {
		fmt.Printf("%s:%d:1 - %s has complexity: %s\n", result.File, result.Line, result.Function, result.Message)
	} else if format == "json" {
		// Add JSON format printing logic here
	} else {
		// Default format if none is specified
		fmt.Printf("%s:%d:1 - %s has complexity: %s\n", result.File, result.Line, result.Function, result.Message)
	}

	// Check if the detailsFormat is true
	if detailsFormat {
		// Get the complexity details for the function
		details := complexity.GetDetail(result)

		// Print the complexity details in a code block
		fmt.Println("```go")
		fmt.Println(details)
		fmt.Println("```")
	}
}

//
//func PrintComplexity(results []*models.LintResult, format string, detailsFormat bool) {
//	var fileCount, funcCount, totalComplexity, maxComplexity, complexLineCount int
//
//	for _, result := range results {
//		if result.File != "" {
//			fileCount++
//		}
//
//		if result.Function != "" {
//			funcCount++
//
//			var complexity int
//			fmt.Sscanf(result.Message, "function has a cognitive complexity of %d", &complexity)
//
//			totalComplexity += complexity
//			if complexity > maxComplexity {
//				maxComplexity = complexity
//			}
//
//			// our assumed threshold is 10,
//			if complexity > 12 {
//				complexLineCount++
//			}
//		}
//	}
//
//	fmt.Printf("%d = files\n%d = functions\n%d = highest complexity\n", fileCount, funcCount, maxComplexity)
//
//	if funcCount == 0 {
//		fmt.Println("No functions were found to calculate the average complexity per function")
//	} else {
//		avgComplexity := float64(totalComplexity) / float64(funcCount)
//
//		fmt.Printf("%.2f = overall average complexity per function\n", avgComplexity)
//	}
//
//	fmt.Printf("%d = complex lines\n\n", complexLineCount)
//
//	// Check if the format is "complexity"
//	if format == "complexity" {
//		// Change the format to "line-number"
//		format = "line-number"
//	}
//
//	for _, result := range results {
//		if result.File != "" && result.Function != "" {
//			if format == "line-number" {
//				fmt.Printf("%s:%d:1 - %s has complexity: %s\n", result.File, result.Line, result.Function, result.Message)
//			} else if format == "json" {
//				// Add JSON format printing logic here
//			} else {
//				// Default format if none is specified
//				fmt.Printf("%s:%d:1 - %s has complexity: %s\n", result.File, result.Line, result.Function, result.Message)
//			}
//			// Check if the detailsFormat is true
//			if detailsFormat {
//				// Get the complexity details for the function
//				details := complexity.GetDetail(result)
//
//				// Print the complexity details in a code block
//				fmt.Println("```go")
//				fmt.Println(details)
//				fmt.Println("```")
//			}
//		}
//	}
//}
