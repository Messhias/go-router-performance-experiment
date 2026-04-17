package features

import (
	"github.com/cucumber/godog"

	"messhias/router-expirement/tests/acceptance_tests"
)

func InitializeScenarios(ctx *godog.ScenarioContext) {
	// init common steps
	acceptance_tests.InitCommon(ctx)

	acceptance_tests.InitOpenAIAcceptanceTests(ctx)
	acceptance_tests.InitRoundRobinLoadBalancing(ctx)
	acceptance_tests.InitProxyTransparency(ctx)
}
