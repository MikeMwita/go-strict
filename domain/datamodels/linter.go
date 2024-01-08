package datamodels

// LintResult represents a linting result for a file or a function

type LintResult struct {
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	Message  string `json:"message,omitempty"`
	Severity string `json:"severity,omitempty"`
	// Fields specific to function linting
	Function string `json:"function,omitempty"`
	// Add other fields as needed
}

// LintConfig represents a linting configuration for the linter
type LintConfig struct {
	Rules     []string // the rules to enable or disable
	Output    string   // the output format (text, json, xml)
	Threshold int
}
