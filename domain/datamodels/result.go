package datamodels

type Result struct {
	File     string
	Line     int
	Message  string
	Severity string // "error" or "warning"
}

//type LintResult struct {
//	// Common fields
//	File     string `json:"file,omitempty"`
//	Line     int    `json:"line,omitempty"`
//	Message  string `json:"message,omitempty"`
//	Severity string `json:"severity,omitempty"`
//
//	// Fields specific to function linting
//	Function string `json:"function,omitempty"`
//	// Add other fields as needed
//}
