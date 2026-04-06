package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MikeMwita/go-strict/models"
)

// FormatReport renders the report in the requested format ("text", "json",
// or "sarif"). It returns the formatted string and any encoding error.
func FormatReport(report *models.LintReport, format string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		return formatJSON(report)
	case "sarif":
		return formatSARIF(report)
	default:
		return formatText(report), nil
	}
}

// FormatSummary returns the summary header for a lint report.
func FormatSummary(report *models.LintReport) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Number of files:     %d\n", report.TotalFiles)
	fmt.Fprintf(&sb, "Number of functions: %d\n", report.TotalFuncs)
	fmt.Fprintf(&sb, "Highest complexity:  %d\n", report.HighestScore)
	if report.TotalFuncs > 0 {
		fmt.Fprintf(&sb, "Average complexity:  %.2f\n", report.AverageScore)
	}
	fmt.Fprintf(&sb, "Complex functions:   %d (threshold: %d)\n", report.ComplexCount, report.Threshold)
	return sb.String()
}

// formatText renders the report in the README-style indented text format.
func formatText(report *models.LintReport) string {
	var sb strings.Builder
	for _, fr := range report.Files {
		for _, fn := range fr.Functions {
			if !fn.Exceeded {
				continue
			}
			fmt.Fprintf(&sb, "%s:%d:%d - %s has complexity: %d\n",
				fn.File, fn.Line, fn.Col, fn.Function, fn.Complexity)

			for _, inc := range fn.Increments {
				indent := strings.Repeat("  ", inc.Nesting)
				fmt.Fprintf(&sb, "%s  + %d (found '%s' at line: %d, complexity = %d)\n",
					indent, inc.Score, inc.Kind, inc.Line, inc.RunningTotal)
			}
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

// formatJSON marshals the full report as indented JSON.
func formatJSON(report *models.LintReport) (string, error) {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// sarif 
type sarifRoot struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string      `json:"name"`
	Version        string      `json:"version"`
	InformationURI string      `json:"informationUri"`
	Rules          []sarifRule `json:"rules"`
}

type sarifRule struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	ShortDescription sarifMessage `json:"shortDescription"`
}

type sarifResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   sarifMessage    `json:"message"`
	Locations []sarifLocation `json:"locations"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysical `json:"physicalLocation"`
}

type sarifPhysical struct {
	ArtifactLocation sarifArtifact `json:"artifactLocation"`
	Region           sarifRegion   `json:"region"`
}

type sarifArtifact struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
}

func formatSARIF(report *models.LintReport) (string, error) {
	root := sarifRoot{
		Version: "2.1.0",
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Runs: []sarifRun{
			{
				Tool: sarifTool{
					Driver: sarifDriver{
						Name:           "go-strict",
						Version:        "1.0.0",
						InformationURI: "https://github.com/MikeMwita/go-strict",
						Rules: []sarifRule{
							{
								ID:               "cognitive-complexity",
								Name:             "CognitiveComplexity",
								ShortDescription: sarifMessage{Text: "Function cognitive complexity exceeds threshold"},
							},
						},
					},
				},
			},
		},
	}

	for _, fr := range report.Files {
		for _, fn := range fr.Functions {
			if !fn.Exceeded {
				continue
			}
			root.Runs[0].Results = append(root.Runs[0].Results, sarifResult{
				RuleID: "cognitive-complexity",
				Level:  "warning",
				Message: sarifMessage{
					Text: fmt.Sprintf(
						"Function %s has cognitive complexity %d (threshold: %d)",
						fn.Function, fn.Complexity, report.Threshold,
					),
				},
				Locations: []sarifLocation{
					{
						PhysicalLocation: sarifPhysical{
							ArtifactLocation: sarifArtifact{URI: fn.File},
							Region:           sarifRegion{StartLine: fn.Line, StartColumn: fn.Col},
						},
					},
				},
			})
		}
	}

	data, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
