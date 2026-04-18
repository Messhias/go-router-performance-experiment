package acceptance_tests

import (
	"messhias/router-expirement/internal/acceptance"

	"github.com/cucumber/godog"
)

func givenUpstreamAAndUpstreamB(harness acceptance.ChatAcceptanceHarness) func() error {
	return func() error {
		return harness.EnsureTwoChatUpstreams()
	}
}

func whenISend4SequentialPost() error {
	return godog.ErrPending
}

func thenUpstreamHandlingOrderShouldBe() error {
	return godog.ErrPending
}

func InitRoundRobinLoadBalancing(ctx *godog.ScenarioContext) {
	ctx.Step(`^I send 4 sequential POST requests to "/v1/chat/completions" with body:$`, whenISend4SequentialPost)
	ctx.Step(`^upstream handling order should be "A,B,A,B"$`, thenUpstreamHandlingOrderShouldBe)
}
