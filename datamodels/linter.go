package datamodels

type LintResult struct {
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	Message  string `json:"message,omitempty"`
	Severity string `json:"severity,omitempty"`
	Function string `json:"function,omitempty"`
}

type LintConfig struct {
	Rules     []string `toml:"rules"`
	Output    string   `toml:"output"`
	Threshold int
}
