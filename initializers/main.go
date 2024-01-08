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

func main() {
	// parse the command-line arguments
	flag.Parse()
	args := flag.Args()

	// check if there are any arguments
	if len(args) == 0 {
		fmt.Println("Usage: go-strict [files or directories]")
		os.Exit(1)
	}

	// load the config file or environment variables
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	// create a complexity service
	complexity := complexity.NewComplexityService()

	// create a linter service
	linter := linter.NewLinterService(config, complexity)

	// Open a file for writing
	file, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Redirect standard output to the file
	os.Stdout = file

	// lint the files or directories
	results, err := linter.LintFiles(args)
	if err != nil {
		fmt.Println("Error linting files:", err)
		os.Exit(1)
	}

	// print the linting results (optional)
	printResults(results, config.Output)
	fmt.Println("Output format:", config.Output)

}

// printResults prints the linting results in the specified format
func printResults(results []*datamodels.LintResult, format string) {
	// Trim leading and trailing spaces from the format
	format = strings.TrimSpace(format)
	switch strings.ToLower(format) {
	case "colored-line-number":
		// print the results with colors and line numbers
	case "line-number":
		// print the results with line numbers
	case "JSON":
		// print the results as a JSON array
		data, err := json.Marshal(results)
		if err != nil {
			fmt.Println("Error encoding results:", err)
			return
		}
		fmt.Println(string(data))
	case "tab":
		// print the results as a tab-separated table
	case "checkstyle":
		// print the results as a checkstyle XML
	case "code-climate":
		// print the results as a code climate JSON
	case "junit-xml":
		// print the results as a JUnit XML
	case "github-actions":
		// print the results as a GitHub actions output
	default:
		// print an error message
		fmt.Println("Invalid output format:", format)
	}
}
