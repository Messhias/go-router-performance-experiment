package acceptance_tests

import (
	"context"
	"messhias/router-expirement/internal/acceptance"
	"messhias/router-expirement/internal/balancer"
	"testing"

	"github.com/cucumber/godog"
)

func givenUpstreamAAndUpstreamB(harness acceptance.ChatAcceptanceHarness) func() error {
	return func() error {
		err := harness.EnsureTwoChatUpstreams()

		if err != nil {
			return err
		}

		a, err := harness.UpstreamAURL()

		if err != nil {
			return err
		}

		b, err := harness.UpstreamBURL()

		if err != nil {
			return err
		}

		roundRobinState.mu.Lock()
		defer roundRobinState.mu.Unlock()

		roundRobinState.upstreamAURL = a
		roundRobinState.upstreamBURL = b
		newBalancer, err := balancer.NewBalancer([]string{roundRobinState.upstreamAURL, roundRobinState.upstreamBURL})

		if err != nil {
			return err
		}

		roundRobinState.bal = newBalancer

		return nil
	}
}

func InitCommon(ctx *godog.ScenarioContext, t *testing.T) {
	harness := acceptance.NewHarness(t)

	ctx.Step(`^router is available$`, givenRouterIsAvailable)
	ctx.Step(`^upstream A and upstream B are configured for chat completions$`, givenUpstreamAAndUpstreamB(harness))
	ctx.Step(`^upstream responds with an OpenAI-compatible chat completion$`, givenUpstreamAAndUpstreamB(harness))

	ctx.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		_ = harness.Close()
		return ctx, nil
	})
}
