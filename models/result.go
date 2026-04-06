package models

// Increment represents a single complexity contribution with its reason,
// score delta, nesting level (for indentation in output), and running total.
type Increment struct {
	Kind         string `json:"kind"`
	Line         int    `json:"line"`
	Col          int    `json:"col"`
	Score        int    `json:"score"`
	Nesting      int    `json:"nesting"`
	RunningTotal int    `json:"running_total"`
}

// FunctionResult holds the complete complexity analysis for a single function.
type FunctionResult struct {
	File       string      `json:"file"`
	Line       int         `json:"line"`
	Col        int         `json:"col"`
	Function   string      `json:"function"`
	Complexity int         `json:"complexity"`
	Increments []Increment `json:"increments,omitempty"`
	Exceeded   bool        `json:"exceeded"`
}

// FileResult aggregates FunctionResults for one source file.
type FileResult struct {
	File      string           `json:"file"`
	Functions []FunctionResult `json:"functions"`
}

// LintReport is the top-level output of a full lint run.
type LintReport struct {
	Files        []FileResult `json:"files"`
	TotalFiles   int          `json:"total_files"`
	TotalFuncs   int          `json:"total_functions"`
	HighestScore int          `json:"highest_complexity"`
	AverageScore float64      `json:"average_complexity"`
	ComplexCount int          `json:"complex_function_count"`
	Threshold    int          `json:"threshold"`
}
