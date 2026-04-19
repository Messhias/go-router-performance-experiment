package acceptance_tests

import "github.com/cucumber/godog"

func givenUpstreamAFailingChatCompletions503() error {
	return godog.ErrPending
}

func thenTheResponseStatusShouldBe(status int) error {
	_ = status
	return godog.ErrPending
}

func thenTheResponseShouldBeValidJSON() error {
	return godog.ErrPending
}

func thenResponseBodyDescribesUpstreamError() error {
	return godog.ErrPending
}

func InitUpstreamFailures(ctx *godog.ScenarioContext) {
	ctx.Step(`^upstream A is failing chat completions with status 503$`, givenUpstreamAFailingChatCompletions503)
	ctx.Step(`^reponse status should be (\d+)$`, thenTheResponseStatusShouldBe)
	ctx.Step(`^reponse should be valid JSON$`, thenTheResponseShouldBeValidJSON)
	ctx.Step(`^reponse body should describe an upstream error$`, thenResponseBodyDescribesUpstreamError)
}
