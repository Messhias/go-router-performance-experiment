package acceptance_tests

import "github.com/cucumber/godog"

func whenClientSendsPOSTWithHeaderAndBody(doc *godog.Table) error {
	_ = doc

	return nil
}

func thenUpstreamHandlingOrderLast3ABABA() error {
	return godog.ErrPending
}

func InitStatelessRouting(ctx *godog.ScenarioContext) {
	ctx.Step(`^following clients send POST "/v1/chat/completions" in order with header "X-Client-Id" and JSON bodies built from:$`, whenClientSendsPOSTWithHeaderAndBody)
	ctx.Step(`^upstream handling order for the last 3 requests should be "A,B,A"$`, thenUpstreamHandlingOrderLast3ABABA)
}
