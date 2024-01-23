package config

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/MikeMwita/go-strict/models"
	"os"
	"path/filepath"
	"strings"
)

type LinterConfig struct {
	Rules  []string `toml:"rules"`
	Output string   `toml:"output"`
}

// LoadConfig loads the linter configuration from the config.toml file
func LoadConfig() (*models.LintConfig, error) {
	configPath := filepath.Join(".", "config.toml")
	if _, err := os.Stat(configPath); err == nil {
		var config models.LintConfig
		if _, err := toml.DecodeFile(configPath, &config); err != nil {
			return nil, err
		}
		return &config, nil
	}
	rules := os.Getenv("LINTER_RULES")
	output := os.Getenv("LINTER_OUTPUT")
	if rules != "" && output != "" {
		config := models.LintConfig{
			Rules:  strings.Split(rules, ","),
			Output: strings.TrimSpace(output),
		}
		return &config, nil
	}

	return nil, errors.New("no config file or environment variables found")
}
