package acceptance_tests

import "github.com/cucumber/godog"

func whenClientSendsPOSTWithHeaderAndBody(client, path, headerName, headerValue string, doc *godog.DocString) error {
	_ = client
	_ = path
	_ = headerName
	_ = headerValue
	_ = doc
	return godog.ErrPending
}

func thenUpstreamHandlingOrderLast3ABABA() error {
	return godog.ErrPending
}

func InitStatelessRouting(ctx *godog.ScenarioContext) {
	ctx.Step(`^client "([^"]*)" sends a POST request to "([^"]*)" with header "([^"]*)" "([^"]*)" and body:$`, whenClientSendsPOSTWithHeaderAndBody)
	ctx.Step(`^upstream handling order for the last 3 requests should be "A,B,A"$`, thenUpstreamHandlingOrderLast3ABABA)
}
