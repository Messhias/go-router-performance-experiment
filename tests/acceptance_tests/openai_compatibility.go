package acceptance_tests

import (
	"bytes"
	"errors"
	"io"
	Dto "messhias/router-expirement/internal/DTO"
	"net/http"
	"net/http/httptest"
	"strconv"

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

	routerTest.srv = httptest.NewServer(serverHandler())

	return nil
}

func givenUpstreamResponds() error {
	return godog.ErrPending
}

func whenPostRequest(doc *godog.DocString) error {
	body := []byte(doc.Content)

	resp, err := http.Post(routerTest.srv.URL+"/v1/chat/completions", "application/json", bytes.NewBuffer(body))

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
	return godog.ErrPending
}

func thenShouldContainOpenAICompatible() error {
	return godog.ErrPending
}

func serverHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/chat/completions" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "read body", http.StatusInternalServerError)
			return
		}

		if _, err := Dto.ParseAndValidateChatRequest(body); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid request"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}
}

func InitOpenAIAcceptanceTests(ctx *godog.ScenarioContext) {
	ctx.Step(`^upstream responds with an OpenAI-compatible chat completion$`, givenUpstreamResponds)
	ctx.Step(`^send a POST request to "/v1/chat/completions" with body:$`, whenPostRequest)
	ctx.Step(`^response status should be 200$`, thenResponseStatus200)
	ctx.Step(`^response status should be 400$`, thenResponseStatus400)
	ctx.Step(`^response should be valid JSON$`, thenShouldValidJson)
	ctx.Step(`^response should contain an OpenAI-compatible chat completion shape$`, thenShouldContainOpenAICompatible)
}
