package exhaustive_test

import (
	"testing"

	"github.com/dannywillems/go-for-systems-engineers/modules/02-sum-types/go/exhaustive"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestExhaustive(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), exhaustive.Analyzer, "a")
}
