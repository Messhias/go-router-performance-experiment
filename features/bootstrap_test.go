package features

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name: "features",
		ScenarioInitializer: func(context *godog.ScenarioContext) {
			InitializeScenarios(context, t)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"."},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatalf("feature tests failed")
	}
}
