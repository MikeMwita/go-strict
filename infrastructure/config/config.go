package config

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/MikeMwita/go-strict/domain/datamodels"
	"os"
	"path/filepath"
	"strings"
)

// LoadConfig loads the lint configuration from either a config file or environment variables
func LoadConfig() (*datamodels.LintConfig, error) {
	// load the config file from the current directory
	configPath := filepath.Join(".", "config.toml")
	if _, err := os.Stat(configPath); err == nil {
		var config datamodels.LintConfig
		if _, err := toml.DecodeFile(configPath, &config); err != nil {
			return nil, err
		}
		return &config, nil
	}

	// loading the config from the environment variables
	rules := os.Getenv("LINTER_RULES")
	output := os.Getenv("LINTER_OUTPUT")
	if rules != "" && output != "" {
		config := datamodels.LintConfig{
			Rules:  strings.Split(rules, ","),
			Output: strings.TrimSpace(output), // Trim spaces
		}
		return &config, nil
	}

	return nil, errors.New("no config file or environment variables found")
}

//
//func LoadConfig() (*datamodels.LintConfig, error) {
//	//  load the config file from the current directory
//	configPath := filepath.Join(".", "config.toml")
//	if _, err := os.Stat(configPath); err == nil {
//		var config datamodels.LintConfig
//		if _, err := toml.DecodeFile(configPath, &config); err != nil {
//			return nil, err
//		}
//		return &config, nil
//	}
//
//	//  loading the config from the environment variables
//
//	rules := os.Getenv("LINTER_RULES")
//	output := os.Getenv("LINTER_OUTPUT")
//	if rules != "" && output != "" {
//		config := datamodels.LintConfig{
//			Rules:  strings.Split(rules, ","),
//			Output: output,
//		}
//		return &config, nil
//	}
//
//	return nil, errors.New("no config file or environment variables found")
//}
