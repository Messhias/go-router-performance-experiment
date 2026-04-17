package features

import (
	"fmt"

	"github.com/cucumber/godog"
)

func theBDDSuiteIsConfigured() error {
	return nil
}

func iRunTheFeatureTests() error {
	return nil
}

func iShouldSeeAFailingScenarioReport() error {
	return fmt.Errorf("intentional failure: validating BDD failure reporting")
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^the BDD suite is configured$`, theBDDSuiteIsConfigured)
	ctx.Step(`^I run the feature tests$`, iRunTheFeatureTests)
	ctx.Step(`^I should see a failing scenario report$`, iShouldSeeAFailingScenarioReport)
}
