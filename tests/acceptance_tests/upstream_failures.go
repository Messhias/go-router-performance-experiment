package acceptance_tests

import (
	"github.com/cucumber/godog"
)

func givenUpstreamAFailingChatCompletions503() error {
	//failSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	if r.URL.Path != config.ChatCompletionsUrl {
	//		http.NotFound(w, r)
	//		return
	//	}
	//}))

	return nil
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
	ctx.Step(`^response status should be (\d+)$`, thenTheResponseStatusShouldBe)
	ctx.Step(`^response should be valid JSON$`, thenTheResponseShouldBeValidJSON)
	ctx.Step(`^response body should describe an upstream error$`, thenResponseBodyDescribesUpstreamError)
}
