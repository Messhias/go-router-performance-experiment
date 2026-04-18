package acceptance_tests

import (
	"context"
	"messhias/router-expirement/internal/acceptance"
	"testing"

	"github.com/cucumber/godog"
)

func givenUpstreamAAndUpstreamB(harness acceptance.ChatAcceptanceHarness) func() error {
	return func() error {
		return harness.EnsureTwoChatUpstreams()
	}
}

func InitCommon(ctx *godog.ScenarioContext, t *testing.T) {
	harness := acceptance.NewHarness(t)

	ctx.Step(`^router is available$`, givenRouterIsAvailable)
	ctx.Step(`^upstream A and upstream B are configured for chat completions$`, givenUpstreamAAndUpstreamB(harness))

	ctx.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		_ = harness.Close()
		return ctx, nil
	})
}
