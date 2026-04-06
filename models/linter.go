package models

// RuleConfig holds enable/disable and threshold settings for a single rule.
type RuleConfig struct {
	Enabled   bool    `toml:"enabled"`
	Weight    float64 `toml:"weight"`
	MaxLines  int     `toml:"max_lines"`
	MinLines  int     `toml:"min_lines"`
	MaxParams int     `toml:"max_params"`
	MaxDepth  int     `toml:"max_depth"`
	Threshold int     `toml:"threshold"`
}

// LintConfig is the top-level configuration structure decoded from config.toml.
type LintConfig struct {
	Threshold    int                   `toml:"threshold"`
	OutputFormat string                `toml:"output_format"`
	OutputFile   string                `toml:"output_file"`
	Rules        map[string]RuleConfig `toml:"rules"`

	Output        string   `toml:"output"`
	MaxComplexity int      `toml:"max_complexity"`
	MaxLineLength int      `toml:"max_line_length"`
	RulesList     []string `toml:"rules_list"`
}

// Rule returns the config for a named rule, or the zero value if not present.
func (c *LintConfig) Rule(name string) RuleConfig {
	if c.Rules == nil {
		return RuleConfig{}
	}
	return c.Rules[name]
}
