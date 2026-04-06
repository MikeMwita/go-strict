package linter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/MikeMwita/go-strict/internal/complexity"
	"github.com/MikeMwita/go-strict/models"
)

// LinterService orchestrates file discovery, parsing, and complexity
// calculation. It processes files in parallel using a worker pool.
type LinterService struct {
	cfg     *models.LintConfig
	workers int
}

// NewLinterService creates a LinterService with the given configuration.
// Pass workers=0 to use runtime.NumCPU().
func NewLinterService(cfg *models.LintConfig, workers int) *LinterService {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	return &LinterService{cfg: cfg, workers: workers}
}

// LintPaths lints all .go files found under the given paths and returns a
// LintReport.
func (ls *LinterService) LintPaths(paths []string) (*models.LintReport, error) {
	files, err := ls.collectFiles(paths)
	if err != nil {
		return nil, err
	}
	return ls.lintFiles(files)
}

// collectFiles walks paths and gathers all .go source files, skipping vendor
// and hidden directories.
func (ls *LinterService) collectFiles(paths []string) ([]string, error) {
	var files []string
	for _, p := range paths {
		err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if info.Name() == "vendor" || strings.HasPrefix(info.Name(), ".") {
					return filepath.SkipDir
				}
				return nil
			}
			if strings.HasSuffix(info.Name(), ".go") &&
				!strings.HasSuffix(info.Name(), "_test.go") {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("walking %s: %w", p, err)
		}
	}
	return files, nil
}

// lintFiles fans file paths out to a worker pool and aggregates results.
func (ls *LinterService) lintFiles(files []string) (*models.LintReport, error) {
	paths := make(chan string, len(files))
	results := make(chan models.FileResult, len(files))
	errs := make(chan error, len(files))

	var wg sync.WaitGroup
	for i := 0; i < ls.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range paths {
				fr, err := ls.lintFile(path)
				if err != nil {
					errs <- err
					continue
				}
				results <- fr
			}
		}()
	}

	for _, f := range files {
		paths <- f
	}
	close(paths)

	wg.Wait()
	close(results)
	close(errs)

	for e := range errs {
		log.Printf("lint error: %v", e)
	}

	report := &models.LintReport{
		TotalFiles: len(files),
		Threshold:  ls.cfg.Threshold,
	}
	for fr := range results {
		report.Files = append(report.Files, fr)
		for _, fn := range fr.Functions {
			report.TotalFuncs++
			report.AverageScore += float64(fn.Complexity)
			if fn.Complexity > report.HighestScore {
				report.HighestScore = fn.Complexity
			}
			if fn.Exceeded {
				report.ComplexCount++
			}
		}
	}
	if report.TotalFuncs > 0 {
		report.AverageScore /= float64(report.TotalFuncs)
	}
	return report, nil
}

// lintFile parses and analyses a single .go file.
func (ls *LinterService) lintFile(path string) (models.FileResult, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return models.FileResult{}, fmt.Errorf("parse %s: %w", path, err)
	}

	fr := models.FileResult{File: path}
	calc := complexity.NewCalculator(fset, f, ls.cfg)

	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil {
			continue
		}
		fr.Functions = append(fr.Functions, calc.Calculate(fd))
	}
	return fr, nil
}
