package utils

import (
	"strings"
	"testing"

	"github.com/MikeMwita/go-strict/models"
)

func TestFormatSummary(t *testing.T) {
	report := &models.LintReport{
		TotalFiles:   3,
		TotalFuncs:   10,
		HighestScore: 15,
		AverageScore: 4.5,
		ComplexCount: 2,
		Threshold:    10,
	}
	out := FormatSummary(report)
	if !strings.Contains(out, "3") {
		t.Error("summary should contain file count")
	}
	if !strings.Contains(out, "15") {
		t.Error("summary should contain highest score")
	}
}

func TestFormatReport_Text(t *testing.T) {
	report := &models.LintReport{
		Threshold: 10,
		Files: []models.FileResult{
			{
				File: "foo.go",
				Functions: []models.FunctionResult{
					{
						File:       "foo.go",
						Line:       1,
						Col:        1,
						Function:   "Heavy",
						Complexity: 12,
						Exceeded:   true,
						Increments: []models.Increment{
							{Kind: "function", Line: 1, Score: 1, Nesting: 0, RunningTotal: 1},
							{Kind: "if", Line: 3, Score: 1, Nesting: 0, RunningTotal: 2},
						},
					},
				},
			},
		},
	}

	out, err := FormatReport(report, "text")
	if err != nil {
		t.Fatalf("FormatReport text: %v", err)
	}
	if !strings.Contains(out, "Heavy has complexity: 12") {
		t.Errorf("expected function name in output, got:\n%s", out)
	}
	if !strings.Contains(out, "found 'if'") {
		t.Errorf("expected if increment in output, got:\n%s", out)
	}
}

func TestFormatReport_JSON(t *testing.T) {
	report := &models.LintReport{Threshold: 10}
	out, err := FormatReport(report, "json")
	if err != nil {
		t.Fatalf("FormatReport json: %v", err)
	}
	if !strings.Contains(out, `"threshold"`) {
		t.Error("JSON output should contain threshold key")
	}
}

func TestFormatReport_SARIF(t *testing.T) {
	report := &models.LintReport{
		Threshold: 10,
		Files: []models.FileResult{
			{
				File: "foo.go",
				Functions: []models.FunctionResult{
					{File: "foo.go", Line: 1, Function: "Heavy", Complexity: 12, Exceeded: true},
				},
			},
		},
	}
	out, err := FormatReport(report, "sarif")
	if err != nil {
		t.Fatalf("FormatReport sarif: %v", err)
	}
	if !strings.Contains(out, `"version": "2.1.0"`) {
		t.Error("SARIF output should contain version 2.1.0")
	}
	if !strings.Contains(out, "Heavy") {
		t.Error("SARIF output should contain function name")
	}
}
