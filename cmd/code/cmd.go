package code

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
	"text": printText, // Custom text output to match the README example
	"json": printJSON,
	"complexity": func(results []*models.LintResult) {
		utils.PrintDetails(results, "complexity", true)
	},
}

func Run() {
	var outputFile string
	flag.StringVar(&outputFile, "o", "", "the output file name")
	var outputFormat string
	flag.StringVar(&outputFormat, "f", "text", "the output format (text, json, complexity)")
	var configPath string
	flag.StringVar(&configPath, "c", "config.toml", "specify the path to the configuration file")
	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "show the version number and exit")
	var showHelp bool
	flag.BoolVar(&showHelp, "h", false, "show the help message and exit")
	flag.Parse()
	args := flag.Args()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Println("Cognitive Complexity Linter v1.0.0")
		os.Exit(0)
	}

	if len(args) == 0 {
		fmt.Println("Usage: cognitive-linter [options] [files or directories]")
		os.Exit(1)
	}

	// Load configuration
	config, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	// Initialize complexity service and linter service
	complexityService := complexity.NewComplexityService()
	linter := linter.NewLinterService(config, complexityService)

	// Handle output file redirection
	var output *os.File = os.Stdout
	if outputFile != "" {
		output, err = os.Create(outputFile)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			os.Exit(1)
		}
		defer output.Close()
	}

	// Convert the relative paths to absolute paths
	var absArgs []string
	for _, arg := range args {
		absArg, err := filepath.Abs(arg)
		if err != nil {
			fmt.Println("Error resolving path:", err)
			os.Exit(1)
		}
		absArgs = append(absArgs, absArg)
	}

	// Lint the files
	results, err := linter.LintFiles(absArgs)
	if err != nil {
		fmt.Println("Error linting files:", err)
		os.Exit(1)
	}

	// Calculate summary statistics
	fileCount := len(absArgs) // Number of files linted
	funcCount := len(results) // Number of functions processed
	highestComplexity := 0
	totalComplexity := 0
	complexLineCount := 0

	for _, result := range results {
		complexity := 0
		fmt.Sscanf(result.Message, "function has a cognitive complexity of %d", &complexity)
		totalComplexity += complexity
		if complexity > highestComplexity {
			highestComplexity = complexity
		}
		if complexity > 12 { // Example threshold for what counts as a "complex" line
			complexLineCount++
		}
	}

	avgComplexity := float64(totalComplexity) / float64(funcCount)

	// Print the summary
	fmt.Fprintf(output, "Number of files: %d\n", fileCount)
	fmt.Fprintf(output, "Number of functions: %d\n", funcCount)
	fmt.Fprintf(output, "Highest complexity: %d\n", highestComplexity)
	fmt.Fprintf(output, "Overall average complexity per function: %.2f\n", avgComplexity)
	fmt.Fprintf(output, "Number of complex lines: %d\n\n", complexLineCount)

	// Print the results in the specified format
	printResults(output, results, outputFormat)
}

func printResults(output *os.File, results []*models.LintResult, format string) {
	format = strings.TrimSpace(format)
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
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Println("Error encoding results:", err)
		return
	}
	fmt.Println(string(data))
}

// printText prints the results in a detailed, structured format like the README example
func printText(results []*models.LintResult) {
	for _, result := range results {
		// Print the basic information about the function
		fmt.Printf("%s:%d:1 - %s\n", result.File, result.Line, result.Function)

		// Extract and print the details of the complexity
		details := strings.Split(result.Message, "Complexity details:\n")
		if len(details) > 1 {
			lines := strings.Split(details[1], "\n")
			for idx, line := range lines {
				trimmedLine := strings.TrimSpace(line)
				if trimmedLine != "" && trimmedLine != ")" {
					// Print each detail correctly formatted
					if strings.HasPrefix(trimmedLine, "+ 1 ") {
						trimmedLine = strings.TrimPrefix(trimmedLine, "+ 1 ")
					}
					if strings.HasSuffix(trimmedLine, ")") {
						fmt.Printf("  + %d (%s)\n", idx+1, trimmedLine)
					} else {
						fmt.Printf("  + %d (%s)\n", idx+1, trimmedLine)
					}
				}
			}
		}

		fmt.Println()
	}
}
