package acceptance_tests

import (
	"context"
	"messhias/router-expirement/internal/balancer"
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

func whenISend4SequentialPost() error {
	return godog.ErrPending
}

func thenUpstreamHandlingOrderShouldBe() error {
	return godog.ErrPending
}

func InitRoundRobinLoadBalancing(ctx *godog.ScenarioContext) {

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		resetRoundRobin()
		return ctx, nil
	})

	ctx.Step(`^I send 4 sequential POST requests to "/v1/chat/completions" with body:$`, whenISend4SequentialPost)
	ctx.Step(`^upstream handling order should be "A,B,A,B"$`, thenUpstreamHandlingOrderShouldBe)
}
