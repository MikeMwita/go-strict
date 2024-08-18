package utils

import (
	"github.com/MikeMwita/go-strict/models"
	"testing"
)

func TestPrintDetails(t *testing.T) {
	type args struct {
		results       []*models.LintResult
		format        string
		detailsFormat bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintDetails(tt.args.results, tt.args.format, tt.args.detailsFormat)
		})
	}
}

func TestPrintSummary(t *testing.T) {
	type args struct {
		fileCount        int
		funcCount        int
		maxComplexity    int
		totalComplexity  int
		funcCountSum     int
		complexLineCount int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintSummary(tt.args.fileCount, tt.args.funcCount, tt.args.maxComplexity, tt.args.totalComplexity, tt.args.funcCountSum, tt.args.complexLineCount)
		})
	}
}
