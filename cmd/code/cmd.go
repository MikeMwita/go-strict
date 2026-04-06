package code

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"github.com/MikeMwita/go-strict/config"
	"github.com/MikeMwita/go-strict/internal/linter"
	"github.com/MikeMwita/go-strict/utils"
)

const version = "1.0.0"

func Run() {
	var (
		outputFile   string
		outputFormat string
		configPath   string
		workers      int
		showVersion  bool
		showHelp     bool
	)

	flag.StringVar(&outputFile, "o", "", "write output to this file (default: stdout)")
	flag.StringVar(&outputFormat, "f", "text", "output format: text | json | sarif")
	flag.StringVar(&configPath, "c", "config.toml", "path to config.toml")
	flag.IntVar(&workers, "p", 0, "parallel workers (default: number of CPUs)")
	flag.BoolVar(&showVersion, "v", false, "print version and exit")
	flag.BoolVar(&showHelp, "h", false, "print help and exit")
	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}
	if showVersion {
		fmt.Printf("go-strict cognitive complexity linter v%s\n", version)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: go-strict [options] <path> [path ...]")
		fmt.Fprintln(os.Stderr, "Run go-strict -h for help.")
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}
	if outputFormat != "" {
		cfg.OutputFormat = outputFormat
	}

	var paths []string
	for _, arg := range args {
		abs, err := filepath.Abs(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "path error: %v\n", err)
			os.Exit(1)
		}
		paths = append(paths, abs)
	}

	svc := linter.NewLinterService(cfg, workers)
	report, err := svc.LintPaths(paths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lint error: %v\n", err)
		os.Exit(1)
	}

	// Choose output destination.
	out := os.Stdout
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "output file error: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		out = f
	}

	// Print summary.
	fmt.Fprint(out, utils.FormatSummary(report))
	fmt.Fprintln(out)

	// Print per-function results.
	formatted, err := utils.FormatReport(report, cfg.OutputFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "format error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprint(out, formatted)

	// Exit with non-zero status if any function exceeded the threshold.
	if report.ComplexCount > 0 {
		os.Exit(1)
	}
}
