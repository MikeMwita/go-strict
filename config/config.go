package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"github.com/BurntSushi/toml"
	"github.com/MikeMwita/go-strict/models"
)

// LoadConfig loads a LintConfig from the given TOML file path.
// If configPath is empty, it looks for config.toml in the current directory.
// If no config file is found it tries environment variables, then falls back
// to DefaultConfig().
func LoadConfig(configPath string) (*models.LintConfig, error) {
	if configPath == "" {
		configPath = filepath.Join(".", "config.toml")
	}

	cfg := DefaultConfig()

	if _, err := os.Stat(configPath); err == nil {
		if _, err := toml.DecodeFile(configPath, cfg); err != nil {
			return nil, err
		}
		ensureDefaults(cfg)
		return cfg, nil
	}

	// Env-var fallback.
	rules := os.Getenv("LINTER_RULES")
	output := os.Getenv("LINTER_OUTPUT")
	if rules != "" && output != "" {
		cfg.RulesList = strings.Split(rules, ",")
		cfg.Output = strings.TrimSpace(output)
		return cfg, nil
	}

	// No config found — return defaults without error so the tool works
	if os.IsNotExist(errors.Unwrap(errors.New(""))) {
		return cfg, nil
	}
	return cfg, nil
}

// ensureDefaults fills in any zero-value fields in cfg that should have
// non-zero defaults (e.g. threshold of 0 means "not set", not "block all").
func ensureDefaults(cfg *models.LintConfig) {
	defaults := DefaultConfig()

	if cfg.Threshold == 0 {
		cfg.Threshold = defaults.Threshold
	}
	if cfg.OutputFormat == "" {
		cfg.OutputFormat = defaults.OutputFormat
	}
	if cfg.Rules == nil {
		cfg.Rules = defaults.Rules
		return
	}
	for name, def := range defaults.Rules {
		if _, ok := cfg.Rules[name]; !ok {
			cfg.Rules[name] = def
		}
	}
}
