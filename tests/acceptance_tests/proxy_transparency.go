package acceptance_tests

import (
	"encoding/json"
	"errors"
	Dto "messhias/router-expirement/internal/DTO"

	"github.com/cucumber/godog"
)

func givenUpstreamAIsConfigured() error {
	return godog.ErrPending
}

func whenSendPostToRequest() error {
	return godog.ErrPending
}

func whenBody() error {
	return godog.ErrPending
}

func thenUpstreamAReceivedJson() error {
	_, err := extractChatResponseAndValidate()

	return err
}

func extractChatResponseAndValidate() (*Dto.ChatCompletionResponseDto, error) {
	if len(routerTest.body) == 0 {
		return nil, errors.New("body is empty")
	}
	var response Dto.ChatCompletionResponseDto
	err := json.Unmarshal(routerTest.body, &response)

	if err != nil {
		return &response, err
	}

	return &response, nil
}
func thenUpstreamAReceivedHeader() error { return godog.ErrPending }

func InitProxyTransparency(ctx *godog.ScenarioContext) {

	ctx.Step(`^upstream A is configured to echo the received request for chat completions$`, givenUpstreamAIsConfigured)
	ctx.Step(`^I send a POST request to "/v1/chat/completions" with headers:$`, whenSendPostToRequest)
	ctx.Step(`^body:$`, whenBody)
	ctx.Step(`^upstream A should have received the same JSON body$`, thenUpstreamAReceivedJson)
	ctx.Step(`^upstream A should have received header "Content-Type" with value "application/json"$`, thenUpstreamAReceivedHeader)
}
