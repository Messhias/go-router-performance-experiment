package acceptance_tests

import (
	"encoding/json"
	"errors"
	"messhias/router-expirement/internal/balancer"
	"messhias/router-expirement/internal/config"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

func givenUpstreamAFailingChatCompletions503() error {
	failSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != config.ChatCompletionsUrl {
			http.NotFound(w, r)
			return
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"message": "upstream server timed out"}`))
	}))

	roundRobinState.mu.Lock()
	defer roundRobinState.mu.Unlock()

	roundRobinState.upstreamAURL = strings.TrimRight(failSrv.URL, "/")

	if roundRobinState.upstreamBURL == "" {
		roundRobinState.upstreamAURL = strings.TrimRight(failSrv.URL, "/")
	}

	newBal, err := balancer.NewBalancer([]string{roundRobinState.upstreamAURL, roundRobinState.upstreamBURL})

	if err != nil {
		return err
	}

	roundRobinState.bal = newBal

	return nil
}

func thenTheResponseStatusShouldBe(status int) error {
	if routerTest.status != status {
		return errors.New("unexpected status code: got " + strconv.Itoa(routerTest.status) + ", want " + strconv.Itoa(status))
	}
	return nil
}

func thenTheResponseShouldBeValidJSON() error {
	if len(routerTest.body) == 0 {
		return errors.New("body is empty")
	}

	var payload map[string]any

	if err := json.Unmarshal(routerTest.body, &payload); err != nil {
		return err
	}

	return nil
}

func thenResponseBodyDescribesUpstreamError() error {
	var payload map[string]any

	if err := json.Unmarshal(routerTest.body, &payload); err != nil {
		return err
	}

	jsonResponse, ok := payload["error"]

	if !ok {
		return errors.New("JSON does not contain error key")
	}

	errorMessage, ok := jsonResponse.(string)

	if !ok || strings.TrimSpace(errorMessage) == "" {
		return errors.New(`"error" should be an non-empty string`)
	}

	lower := strings.ToLower(errorMessage)

	if !containsAny(lower, "upstream", "timeout", "unavailable") {
		return errors.New(`"upstream" should contain "timeout" key`)
	}

	return nil
}

func containsAny(word string, terms ...string) bool {

	for _, term := range terms {
		if strings.Contains(word, term) {
			return true
		}
	}

	return false
}

func InitUpstreamFailures(ctx *godog.ScenarioContext) {
	ctx.Step(`^upstream A is failing chat completions with status 503$`, givenUpstreamAFailingChatCompletions503)
	ctx.Step(`^response status should be (\d+)$`, thenTheResponseStatusShouldBe)
	ctx.Step(`^response should be valid JSON$`, thenTheResponseShouldBeValidJSON)
	ctx.Step(`^response body should describe an upstream error$`, thenResponseBodyDescribesUpstreamError)
}
