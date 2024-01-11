package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/MikeMwita/go-strict/domain/datamodels"
	"github.com/MikeMwita/go-strict/infrastructure/config"
	"github.com/MikeMwita/go-strict/services/complexity"
	"github.com/MikeMwita/go-strict/services/linter"
	"os"
	"strings"
)

// define a map of output format names and printing functions
var outputFormats = map[string]func([]*datamodels.LintResult){
	"line-number": printLineNumbers,
	"json":        printJSON,
	"complexity":  printComplexity,
}

func main() {
	// flag for the output file name
	var outputFile string
	flag.StringVar(&outputFile, "o", "", "the output file name")

	//  flag for the output format
	var outputFormat string
	flag.StringVar(&outputFormat, "f", "", "the output format")
	// parse the command-line arguments
	flag.Parse()
	args := flag.Args()

	// check if there are any arguments
	if len(args) == 0 {
		fmt.Println("Usage: go-strict [files or directories]")
		os.Exit(1)
	}

	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	//  complexity service
	complexity := complexity.NewComplexityService()

	//  linter service
	linter := linter.NewLinterService(config, complexity)

	// check if the output file name is given
	if outputFile != "" {
		// create the output file
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			os.Exit(1)
		}
		defer file.Close()

		// redirect the standard output to the file
		os.Stdout = file
	}

	// lint the files or directories
	results, err := linter.LintFiles(args)
	if err != nil {
		fmt.Println("Error linting files:", err)
		os.Exit(1)
	}

	// print the linting results (optional)
	printResults(results, outputFormat)
	fmt.Println("Output format:", outputFormat)
}

func printResults(results []*datamodels.LintResult, format string) {
	// Trim leading and trailing spaces from the format
	format = strings.TrimSpace(format)

	// check if the format is empty
	if format == "" {
		// use the default output format from the config file or environment variable
		//format = config.LintConfig.Output
		format = "JSON"

	}

	// convert the format to lower case
	format = strings.ToLower(format)

	// look up the printing function for the format
	printFunc, ok := outputFormats[format]
	if !ok {
		// print an error message if the format is not found
		fmt.Println("Invalid output format:", format)
		return
	}
	printFunc(results)
}

// printComplexity prints the results in the complexity format
func printComplexity(results []*datamodels.LintResult) {
	// create some variables to store the statistics
	var fileCount, funcCount, totalComplexity, maxComplexity, complexLineCount int

	// iterate over the results
	for _, result := range results {
		// check if the result has a file name
		if result.File != "" {
			// increment the file count
			fileCount++
		}

		// check if the result has a function name
		if result.Function != "" {
			// increment the function count
			funcCount++

			var complexity int
			fmt.Sscanf(result.Message, "function has a cognitive complexity of %d", &complexity)

			// update the total and maximum complexity
			totalComplexity += complexity
			if complexity > maxComplexity {
				maxComplexity = complexity
			}

			// check if the complexity exceeds the threshold
			// our assumed threshold is 10,
			if complexity > 10 {
				// increment the complex line count
				complexLineCount++
			}
		}
	}

	// print the statistics
	fmt.Printf("%d = files\n%d = functions\n%d = highest complexity\n", fileCount, funcCount, maxComplexity)

	// check if the function count is 0
	if funcCount == 0 {
		// print a message indicating that no functions were found
		fmt.Println("No functions were found to calculate the average complexity per function")
	} else {
		// calculate the average complexity
		avgComplexity := float64(totalComplexity) / float64(funcCount)

		// print the average complexity
		fmt.Printf("%.2f = overall average complexity per function\n", avgComplexity)
	}

	// print the complex line count
	fmt.Printf("%d = complex lines\n\n", complexLineCount)

	// print the results
	for _, result := range results {
		// check if the result has a file name and a function name
		if result.File != "" && result.Function != "" {
			// print the result with the file name, line number, function name, and message
			fmt.Printf("%s:%d:1 - %s has complexity: %s\n", result.File, result.Line, result.Function, result.Message)
		}
	}
}

// printJSON prints the results as JSON
func printJSON(results []*datamodels.LintResult) {
	// encode the results as JSON
	data, err := json.Marshal(results)
	if err != nil {
		fmt.Println("Error encoding results:", err)
		return
	}
	// print the JSON data
	fmt.Println(string(data))
}

// printLineNumbers prints the results with line numbers
func printLineNumbers(results []*datamodels.LintResult) {
	// iterate over the results
	for _, result := range results {
		// check if the result has a line number
		if result.Line > 0 {
			// print the result with the line number
			fmt.Printf("%s:%d: %s\n", result.File, result.Line, result.Message)
		} else {
			// print the result without the line number
			fmt.Printf("%s: %s\n", result.File, result.Message)
		}
	}
}
