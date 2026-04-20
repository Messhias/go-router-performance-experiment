package acceptance_tests

import (
	"bytes"
	"errors"
	"io"
	"messhias/router-expirement/internal/config"
	"messhias/router-expirement/internal/proxy"
	"messhias/router-expirement/internal/router"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

var routerTest struct {
	status int
	body   []byte
	srv    *httptest.Server
}

func givenRouterIsAvailable() error {
	if routerTest.srv != nil {
		routerTest.srv.Close()
	}

	hooks := &proxy.Hooks{
		OnUpstream2xx: func(chosen string) {
			roundRobinState.mu.Lock()
			a := strings.TrimRight(roundRobinState.upstreamAURL, "/")
			b := strings.TrimRight(roundRobinState.upstreamBURL, "/")
			roundRobinState.mu.Unlock()

			switch strings.TrimRight(chosen, "/") {
			case a:
				appendRoundRobinHandlingLetter("A")
			case b:
				appendRoundRobinHandlingLetter("B")
			}
		},
	}

	routerTest.srv = httptest.NewServer(router.NewEngine(&delegatingBalancer{}, hooks))
	return nil
}

func whenPostRequest(doc *godog.DocString) error {
	body := []byte(doc.Content)

	resp, err := http.Post(routerTest.srv.URL+config.ChatCompletionsUrl, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	routerTest.status = resp.StatusCode
	routerTest.body, err = io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	return nil
}

func thenResponseStatus200() error {
	if routerTest.status != http.StatusOK {
		return errors.New("unexpected status code: " + strconv.Itoa(routerTest.status))
	}

	return nil
}

func thenResponseStatus400() error {
	if routerTest.status != http.StatusBadRequest {
		return errors.New("unexpected status code: " + strconv.Itoa(routerTest.status))
	}

	return nil
}

func thenShouldValidJson() error {
	_, err := extractChatResponseAndValidate()

	if err != nil {
		return err
	}

	return nil
}

func thenShouldContainOpenAICompatible() error {
	chatResponse, err := extractChatResponseAndValidate()

	if err != nil {
		return err
	}

	if chatResponse.Object != "chat.completion" {
		return errors.New("expected object to be 'chat.completion'")
	}

	if len(chatResponse.Choices) == 0 {
		return errors.New("expected choices to contain at least one choice")
	}

	return nil
}

func InitOpenAIAcceptanceTests(ctx *godog.ScenarioContext) {
	ctx.Step(`^send a POST request to "/v1/chat/completions" with body:$`, whenPostRequest)
	ctx.Step(`^response status should be 200$`, thenResponseStatus200)
	ctx.Step(`^response status should be 400$`, thenResponseStatus400)
	ctx.Step(`^response should be valid JSON$`, thenShouldValidJson)
	ctx.Step(`^response should contain an OpenAI-compatible chat completion shape$`, thenShouldContainOpenAICompatible)
}
