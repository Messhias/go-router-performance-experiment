package acceptance_tests

import "github.com/cucumber/godog"

func givenUpstreamAIsConfigured() error {
	return godog.ErrPending
}

func whenSendPostToRequest() error {
	return godog.ErrPending
}

func whenBody() error {
	return godog.ErrPending
}

func thenUpstreamAReceivedJson() error   { return godog.ErrPending }
func thenUpstreamAReceivedHeader() error { return godog.ErrPending }

func InitProxyTransparency(ctx *godog.ScenarioContext) {

	ctx.Step(`^router is available$`, givenRouterIsAvailable)
	ctx.Step(`^upstream A is configured to echo the received request for chat completions$`, givenUpstreamAIsConfigured)
	ctx.Step(`^I send a POST request to "/v1/chat/completions" with headers:$`, whenSendPostToRequest)
	ctx.Step(`^body:$`, whenBody)
	ctx.Step(`^upstream A should have received the same JSON body$`, thenUpstreamAReceivedJson)
	ctx.Step(`^upstream A should have received header "Content-Type" with value "application/json"$`, thenUpstreamAReceivedHeader)
}
