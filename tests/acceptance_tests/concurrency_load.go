package acceptance_tests

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"messhias/router-expirement/internal/config"
	"net/http"
	"strings"
	"sync"

	"github.com/cucumber/godog"
)

func whenSendNParallelPOSTRequests(n int, doc *godog.DocString) error {
	if routerTest.srv == nil {
		return errors.New("router is nil")
	}

	body := []byte(strings.TrimSpace(doc.Content))
	url := routerTest.srv.URL + config.ChatCompletionsUrl

	roundRobinState.mu.Lock()
	roundRobinState.handlingOrder = nil
	roundRobinState.mu.Unlock()

	var wg sync.WaitGroup
	var errMu sync.Mutex
	var firstErr error
	fail := func(format string, args ...any) {
		errMu.Lock()
		if firstErr == nil {
			firstErr = fmt.Errorf(format, args...)
		}
		errMu.Unlock()
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Post(url, "application/json", bytes.NewReader(body))
			if err != nil {
				fail("post: %w", err)
				return
			}
			defer func() { _ = resp.Body.Close() }()

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				fail("read body: %w", err)
				return
			}
			if resp.StatusCode != http.StatusOK {
				fail("status %d: %s", resp.StatusCode, string(respBody))
			}
		}()
	}
	wg.Wait()

	errMu.Lock()
	e := firstErr
	errMu.Unlock()
	if e != nil {
		return e
	}

	roundRobinState.mu.Lock()
	got := len(roundRobinState.handlingOrder)
	roundRobinState.mu.Unlock()
	if got != n {
		return fmt.Errorf("expected %d successful upstream selections, got %d", n, got)
	}

	return nil
}

func thenUpstreamABEachBetweenPercent(minPct, maxPct int) error {
	roundRobinState.mu.Lock()
	handlers := append([]string(nil), roundRobinState.handlingOrder...)
	roundRobinState.mu.Unlock()

	var countA, countB int
	for _, h := range handlers {
		switch h {
		case "A":
			countA++
		case "B":
			countB++
		default:
			return fmt.Errorf("unexpected upstream label %q", h)
		}
	}

	total := countA + countB
	if total == 0 {
		return errors.New("no upstream handling recorded")
	}

	pctA := 100.0 * float64(countA) / float64(total)
	pctB := 100.0 * float64(countB) / float64(total)
	minF := float64(minPct)
	maxF := float64(maxPct)

	if pctA < minF || pctA > maxF {
		return fmt.Errorf("upstream A share %.1f%% not in [%d,%d]%% (A=%d B=%d)", pctA, minPct, maxPct, countA, countB)
	}
	if pctB < minF || pctB > maxF {
		return fmt.Errorf("upstream B share %.1f%% not in [%d,%d]%% (A=%d B=%d)", pctB, minPct, maxPct, countA, countB)
	}

	return nil
}

// InitConcurrencyLoad registers BDD steps for parallel load scenarios.
func InitConcurrencyLoad(ctx *godog.ScenarioContext) {
	ctx.Step(`^send (\d+) parallel POST requests to "/v1/chat/completions" with body:$`, whenSendNParallelPOSTRequests)
	ctx.Step(`^upstream A and upstream B should each handle between (\d+) and (\d+) percent of requests$`, thenUpstreamABEachBetweenPercent)
}
