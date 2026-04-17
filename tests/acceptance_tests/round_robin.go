package acceptance_tests

import "github.com/cucumber/godog"

func givenUpstreamAAndUpstreamB() error {
	return godog.ErrPending
}

func whenISend4SequentialPost() error {
	return godog.ErrPending
}

func thenUpstreamHandlingOrderShouldBe() error {
	return godog.ErrPending
}

func InitRoundRobinLoadBalancing(ctx *godog.ScenarioContext) {

	ctx.Step(`^router is available$`, givenRouterIsAvailable)
	ctx.Step(`^upstream A and upstream B are configured for chat completions$`, givenUpstreamAAndUpstreamB)
	ctx.Step(`^I send 4 sequential POST requests to "/v1/chat/completions" with body:$`, whenISend4SequentialPost)
	ctx.Step(`^upstream handling order should be "A,B,A,B"$`, thenUpstreamHandlingOrderShouldBe)
}
