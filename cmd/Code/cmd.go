package Code

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/MikeMwita/go-strict/config"
	"github.com/MikeMwita/go-strict/internal/linter"
	"github.com/MikeMwita/go-strict/models"
	"github.com/MikeMwita/go-strict/services/complexity"
	"github.com/MikeMwita/go-strict/utils"
	"os"
	"path/filepath"
	"strings"
)

var outputFormats = map[string]func([]*models.LintResult){
	"line-number": printLineNumbers,
	"json":        printJSON,
	"complexity":  printComplexity,
}

func Run() {
	var outputFile string
	flag.StringVar(&outputFile, "o", "", "the output file name")
	var outputFormat string
	flag.StringVar(&outputFormat, "f", "", "the output format")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Usage: go-strict [files or directories]")
		os.Exit(1)
	}

	//  original working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		os.Exit(1)
	}

	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	complexity := complexity.NewComplexityService()

	linter := linter.NewLinterService(config, complexity)

	if outputFile != "" {
		//  output file
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			os.Exit(1)
		}
		defer file.Close()

		os.Stdout = file
	}

	// Convert the relative paths to absolute paths
	var absArgs []string
	for _, arg := range args {
		absArg := filepath.Join(wd, arg)
		absArgs = append(absArgs, absArg)
	}

	results, err := linter.LintFiles(absArgs)
	if err != nil {
		fmt.Println("Error linting files:", err)
		os.Exit(1)
	}

	printResults(results, outputFormat)
	fmt.Println("Output format:", outputFormat)
}

func printResults(results []*models.LintResult, format string) {
	format = strings.TrimSpace(format)
	if format == "" {
		format = "JSON"

	}

	format = strings.ToLower(format)

	printFunc, ok := outputFormats[format]
	if !ok {
		fmt.Println("Invalid output format:", format)
		return
	}
	printFunc(results)
}

// printJSON prints the results as JSON

func printJSON(results []*models.LintResult) {
	data, err := json.Marshal(results)
	if err != nil {
		fmt.Println("Error encoding results:", err)
		return
	}
	fmt.Println(string(data))
}

// printLineNumbers prints the results with line numbers
func printLineNumbers(results []*models.LintResult) {
	for _, result := range results {
		if result.Line > 0 {
			fmt.Printf("%s:%d: %s\n", result.File, result.Line, result.Message)
		} else {
			fmt.Printf("%s: %s\n", result.File, result.Message)
		}
	}
}

func printComplexity(results []*models.LintResult) {
	utils.PrintComplexity(results, "complexity", true)
}
