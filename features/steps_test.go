package features

import (
	"testing"

	"github.com/cucumber/godog"

	"messhias/router-expirement/tests/acceptance_tests"
)

func InitializeScenarios(ctx *godog.ScenarioContext, t *testing.T) {
	// init common steps
	acceptance_tests.InitCommon(ctx, t)

	acceptance_tests.InitOpenAIAcceptanceTests(ctx)
	acceptance_tests.InitRoundRobinLoadBalancing(ctx)
	acceptance_tests.InitProxyTransparency(ctx)
	acceptance_tests.InitStatelessRouting(ctx)
	acceptance_tests.InitUpstreamFailures(ctx)
}
