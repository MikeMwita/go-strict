package datamodels

// LintResult represents a linting result for a file or a function

type LintResult struct {
	File     string // the file name
	Line     int    // the line number
	Message  string // the linting message
	Severity string // the severity level (info, warning, error)
}

// LintConfig represents a linting configuration for the linter
type LintConfig struct {
	Rules     []string // the rules to enable or disable
	Output    string   // the output format (text, json, xml)
	Threshold int
}
