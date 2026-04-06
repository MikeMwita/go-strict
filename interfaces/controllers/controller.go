package controllers

import (
	"net/http"

	"github.com/MikeMwita/go-strict/internal/linter"
	"github.com/gin-gonic/gin"
)

type LintController struct {
	linterService *linter.LinterService
}

func NewLintController(linterService *linter.LinterService) *LintController {
	return &LintController{linterService: linterService}
}

// LintPaths lints the given file or directory paths.
func (lc *LintController) LintPaths(c *gin.Context) {
	paths := c.QueryArray("paths")
	if len(paths) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no paths given"})
		return
	}

	report, err := lc.linterService.LintPaths(paths)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}
