package acceptance_tests

import "github.com/cucumber/godog"

func whenClientSendsPOSTWithHeaderAndBody(table *godog.Table) error {

	for _, row := range table.Rows[1:] {
		if len(row.Cells) < 2 {
			continue
		}

	}

	return nil
}

func thenUpstreamHandlingOrderLast3ABABA() error {
	return godog.ErrPending
}

func InitStatelessRouting(ctx *godog.ScenarioContext) {
	ctx.Step(`^following clients send POST "/v1/chat/completions" in order with header "X-Client-Id" and JSON bodies built from:$`, whenClientSendsPOSTWithHeaderAndBody)
	ctx.Step(`^upstream handling order for the last 3 requests should be "A,B,A"$`, thenUpstreamHandlingOrderLast3ABABA)
}
