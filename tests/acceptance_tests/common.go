package acceptance_tests

import (
	"testing"

	"github.com/cucumber/godog"
)

func InitCommon(ctx *godog.ScenarioContext, t *testing.T) {

	ctx.Step(`^router is available$`, givenRouterIsAvailable)
	ctx.Step(`^upstream A and upstream B are configured for chat completions$`, givenUpstreamAAndUpstreamB)
}
