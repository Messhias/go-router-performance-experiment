package acceptance_tests

import (
	"bytes"
	"errors"
	"io"
	Dto "messhias/router-expirement/internal/DTO"
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

	routerTest.srv = httptest.NewServer(serverHandler())

	return nil
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

		roundRobinState.mu.Lock()
		bal := roundRobinState.bal
		urlA := roundRobinState.upstreamAURL
		urlB := roundRobinState.upstreamBURL
		roundRobinState.mu.Unlock()

		if bal == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
			return
		}

		target, err := bal.Next()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		upURL := strings.TrimRight(target, "/") + "/v1/chat/completions"

		req, err := http.NewRequestWithContext(r.Context(), http.MethodPost, upURL, bytes.NewReader(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if ct := r.Header.Get("Content-Type"); ct != "" {
			req.Header.Set("Content-Type", ct)
		} else {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		upBody, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			t := strings.TrimRight(target, "/")
			switch t {
			case strings.TrimRight(urlA, "/"):
				appendRoundRobinHandlingLetter("A")
			case strings.TrimRight(urlB, "/"):
				appendRoundRobinHandlingLetter("B")
			}
		}

		if ct := resp.Header.Get("Content-Type"); ct != "" {
			w.Header().Set("Content-Type", ct)
		}
		w.WriteHeader(resp.StatusCode)
		_, _ = w.Write(upBody)
	}
}

func InitOpenAIAcceptanceTests(ctx *godog.ScenarioContext) {
	ctx.Step(`^send a POST request to "/v1/chat/completions" with body:$`, whenPostRequest)
	ctx.Step(`^response status should be 200$`, thenResponseStatus200)
	ctx.Step(`^response status should be 400$`, thenResponseStatus400)
	ctx.Step(`^response should be valid JSON$`, thenShouldValidJson)
	ctx.Step(`^response should contain an OpenAI-compatible chat completion shape$`, thenShouldContainOpenAICompatible)
}
