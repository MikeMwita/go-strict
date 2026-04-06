package config

import "github.com/MikeMwita/go-strict/models"

// DefaultConfig returns a fully-populated LintConfig with sensible defaults
// for every rule. It is used as the base when a config file is absent or
// when specific keys are missing from the file.
func DefaultConfig() *models.LintConfig {
	return &models.LintConfig{
		Threshold:    10,
		OutputFormat: "text",
		Rules: map[string]models.RuleConfig{
			"cognitive_complexity": {
				Enabled:   true,
				Threshold: 10,
			},
			"function_length": {
				Enabled:  true,
				Weight:   1,
				MaxLines: 50,
			},
			"parameter_count": {
				Enabled:   true,
				Weight:    1,
				MaxParams: 5,
			},
			"missing_comment": {
				Enabled:  true,
				Weight:   1,
				MinLines: 10,
			},
			"nesting_depth": {
				Enabled:  true,
				Weight:   1,
				MaxDepth: 4,
			},
			"logical_operators": {
				Enabled: true,
				Weight:  1,
			},
			"generics": {
				Enabled: true,
				Weight:  1,
			},
			"block_length": {
				Enabled:  true,
				Weight:   1,
				MinLines: 10,
			},
		},
	}
}
