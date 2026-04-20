package acceptance_tests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"messhias/router-expirement/internal/balancer"
	"messhias/router-expirement/internal/config"
	"net/http"
	"strings"
	"sync"

	"github.com/cucumber/godog"
)

var roundRobinState struct {
	mu sync.Mutex

	upstreamAURL string
	upstreamBURL string

	handlingOrder []string

	bal balancer.RoundRobin
}

func resetRoundRobin() {
	roundRobinState.mu.Lock()
	defer roundRobinState.mu.Unlock()

	roundRobinState.upstreamAURL = ""
	roundRobinState.upstreamBURL = ""

	roundRobinState.bal = nil
	roundRobinState.handlingOrder = nil
}

func appendRoundRobinHandlingLetter(letter string) {
	roundRobinState.mu.Lock()
	defer roundRobinState.mu.Unlock()

	roundRobinState.handlingOrder = append(roundRobinState.handlingOrder, letter)
}

func whenISend4SequentialPost(doc *godog.DocString) error {
	if routerTest.srv == nil {
		return errors.New("router is nil")
	}

	body := []byte(strings.TrimSpace(doc.Content))
	url := routerTest.srv.URL + config.ChatCompletionsUrl

	roundRobinState.mu.Lock()
	roundRobinState.handlingOrder = nil
	roundRobinState.mu.Unlock()

	for i := 0; i < 4; i++ {
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

		if err != nil {
			return fmt.Errorf("could not create request: %w in interaction %d", err, i)
		}

		response, readErr := io.ReadAll(resp.Body)
		err = resp.Body.Close()

		if err != nil {
			return fmt.Errorf("could not close request: %w in %d", err, i)
		}

		if readErr != nil {
			return fmt.Errorf("could not close request: %w in %d", readErr, i)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("request %d: status %d body %s", i, resp.StatusCode, string(response))
		}
	}

	return nil
}

func thenUpstreamHandlingOrderShouldBe() error {
	want := []string{"A", "B", "A", "B"}
	roundRobinState.mu.Lock()
	defer roundRobinState.mu.Unlock()
	got := append([]string(nil), roundRobinState.handlingOrder...)
	if len(got) != len(want) {
		return fmt.Errorf("handling order length %d, got %v, want %v", len(got), got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			return fmt.Errorf("handling order at %d: got %q want %q (full got %v)", i, got[i], want[i], got)
		}
	}
	return nil
}

func InitRoundRobinLoadBalancing(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		resetRoundRobin()
		return ctx, nil
	})

	ctx.Step(`^I send 4 sequential POST requests to "/v1/chat/completions" with body:$`, whenISend4SequentialPost)
	ctx.Step(`^upstream handling order should be "A,B,A,B"$`, thenUpstreamHandlingOrderShouldBe)
}
