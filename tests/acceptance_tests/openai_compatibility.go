package acceptance_tests

import "github.com/cucumber/godog"

func givenRouterIsAvailable() error {
	return godog.ErrPending
}

func givenUpstreamResponds() error {
	return godog.ErrPending
}

func whenPostRequest() error {
	return godog.ErrPending
}

func thenResponseStatus200() error {
	return godog.ErrPending
}

func thenShouldValidJson() error {
	return godog.ErrPending
}

func thenShouldContainOpenAICompatible() error {
	return godog.ErrPending
}

func InitOpenAIAcceptanceTests(ctx *godog.ScenarioContext) {
	ctx.Step(`^upstream responds with an OpenAI-compatible chat completion$`, givenUpstreamResponds)
	ctx.Step(`^send a POST request to "/v1/chat/completions" with body:$`, whenPostRequest)
	ctx.Step(`^response status should be 200$`, thenResponseStatus200)
	ctx.Step(`^response should be valid JSON$`, thenShouldValidJson)
	ctx.Step(`^response should contain an OpenAI-compatible chat completion shape$`, thenShouldContainOpenAICompatible)
}
