package linter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MikeMwita/go-strict/config"
)

func newTestService() *LinterService {
	cfg := config.DefaultConfig()
	cfg.Threshold = 5
	return NewLinterService(cfg, 1)
}

func TestLintPaths_SimpleFunction(t *testing.T) {
	src := `package p
func Add(x, y int) int { return x + y }
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "add.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	report, err := newTestService().LintPaths([]string{filepath.Join(dir, "add.go")})
	if err != nil {
		t.Fatalf("LintPaths: %v", err)
	}
	if report.ComplexCount != 0 {
		t.Errorf("simple function should not exceed threshold, got ComplexCount=%d", report.ComplexCount)
	}
}

func TestLintPaths_ComplexFunction(t *testing.T) {
	src := `package p

// Foo does something complex.
func Foo(a, b, c bool) bool {
	if a {
		if b {
			if c {
				return true
			}
		}
	}
	return false
}
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "foo.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	report, err := newTestService().LintPaths([]string{dir})
	if err != nil {
		t.Fatalf("LintPaths: %v", err)
	}
	if report.TotalFiles != 1 {
		t.Errorf("want 1 file, got %d", report.TotalFiles)
	}
	if report.ComplexCount == 0 {
		t.Error("Foo should exceed threshold of 5")
	}
}

func TestLintPaths_SkipsVendor(t *testing.T) {
	dir := t.TempDir()
	vendor := filepath.Join(dir, "vendor")
	if err := os.MkdirAll(vendor, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vendor, "dep.go"), []byte("package dep\nfunc Dep() {}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\nfunc main() {}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	report, err := newTestService().LintPaths([]string{dir})
	if err != nil {
		t.Fatalf("LintPaths: %v", err)
	}
	if report.TotalFiles != 1 {
		t.Errorf("vendor should be skipped: want 1 file, got %d", report.TotalFiles)
	}
}

func TestLintPaths_ReturnsReport(t *testing.T) {
	src := `package p
func A() {}
func B() {}
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "ab.go"), []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	report, err := newTestService().LintPaths([]string{dir})
	if err != nil {
		t.Fatalf("LintPaths: %v", err)
	}
	if report.TotalFuncs != 2 {
		t.Errorf("want 2 functions, got %d", report.TotalFuncs)
	}
}
