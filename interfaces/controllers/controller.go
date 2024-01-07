package controllers

import (
	"github.com/MikeMwita/go-strict/services/linter"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LintController is a controller that handles the linting requests
type LintController struct {
	linterService linter.LinterService // the linter service
}

// NewLintController creates a new LintController
func NewLintController(linterService linter.LinterService) *LintController {
	return &LintController{
		linterService: linterService,
	}
}

// LintFiles is a handler that lints the given files or directories
func (lc *LintController) LintFiles(c *gin.Context) {
	// get the files or directories from the query parameters
	files := c.QueryArray("files")

	// check if the files or directories are given
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no files or directories given",
		})
		return
	}

	// lint the files or directories using the linter service
	results, err := lc.linterService.LintFiles(files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// return the linting results as JSON
	c.JSON(http.StatusOK, results)
}

// LintFunctions is a handler that lints the given functions
func (lc *LintController) LintFunctions(c *gin.Context) {
	// get the functions from the query parameters
	functions := c.QueryArray("functions")

	// check if the functions are given
	if len(functions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "no functions given",
		})
		return
	}

	// lint the functions using the linter service
	results, err := lc.linterService.LintFunctions(functions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// return the linting results as JSON
	c.JSON(http.StatusOK, results)
}
