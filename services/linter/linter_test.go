package linter_test

import (
	"github.com/MikeMwita/go-strict/services/linter"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestNodeCount(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, linter.NodeCount, "a")
}
