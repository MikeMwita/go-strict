package controllers

import (
	"github.com/MikeMwita/go-strict/internal/linter"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LintController struct {
	linterService linter.LinterService // the linter service
}

func (lc *LintController) LintFiles(c *gin.Context) {
	files := c.QueryArray("files")
	// check if the files or directories are given
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no files or directories given",
		})
		return
	}

	results, err := lc.linterService.LintFiles(files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, results)
}

// LintFunctions is a handler that lints the given functions
func (lc *LintController) LintFunctions(c *gin.Context) {
	functions := c.QueryArray("functions")

	if len(functions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no functions given",
		})
		return
	}

	results, err := lc.linterService.LintFunctions(functions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, results)
}

func NewLintController(linterService linter.LinterService) *LintController {
	return &LintController{
		linterService: linterService,
	}
}
