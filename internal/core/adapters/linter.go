package adapters

import "github.com/MikeMwita/go-strict/models"

// Linter is the public interface for the linting service.
type Linter interface {
	LintPaths(paths []string) (*models.LintReport, error)
}
