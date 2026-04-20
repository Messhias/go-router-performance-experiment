package acceptance_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	Dto "messhias/router-expirement/internal/DTO"
	"messhias/router-expirement/internal/balancer"
	"messhias/router-expirement/internal/config"
	"messhias/router-expirement/internal/proxy"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

func givenUpstreamEchoForChatCompletions(th proxy.Transparency) error {
	mu := th.GetMu()
	mu.Lock()
	if th.GetEchoSrv() != nil {
		th.GetEchoSrv().Close()
		th.SetEchoSrv(nil)
	}
	th.SetLastUpstreamBody(nil)
	th.SetLastUpstreamCT("")
	th.SetLastClientBody(nil)
	th.SetPendingHeaders(nil)
	mu.Unlock()

	echo := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != config.ChatCompletionsUrl {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mu.Lock()
		th.SetLastUpstreamBody(append([]byte(nil), body...))
		th.SetLastUpstreamCT(r.Header.Get("Content-Type"))
		mu.Unlock()

		resp := Dto.ChatCompletionResponseDto{
			ID:      "chatcmpl-echo",
			Object:  "chat.completion",
			Created: 1,
			Model:   "echo",
			Choices: []Dto.ChoiceDto{{
				Index:        0,
				FinishReason: "stop",
				Message:      Dto.MessageDto{Role: "assistant", Content: "ok"},
			}},
		}
		payload, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
	}))

	mu.Lock()
	th.SetEchoSrv(echo)
	mu.Unlock()

	roundRobinState.mu.Lock()
	defer roundRobinState.mu.Unlock()
	base := strings.TrimRight(echo.URL, "/")
	roundRobinState.upstreamAURL = base
	roundRobinState.upstreamBURL = base
	bal, err := balancer.NewBalancer([]string{echo.URL, echo.URL})
	if err != nil {
		echo.Close()
		return err
	}
	roundRobinState.bal = bal
	return nil
}

func whenSendPostToChatCompletionsWithHeaders(th proxy.Transparency, table *godog.Table) error {
	if table == nil || len(table.Rows) < 2 {
		return errors.New("headers table must have a header row and at least one data row")
	}
	h := http.Header{}
	for _, row := range table.Rows[1:] {
		if len(row.Cells) < 2 {
			continue
		}
		name := strings.TrimSpace(row.Cells[0].Value)
		val := strings.TrimSpace(row.Cells[1].Value)
		h.Set(name, val)
	}
	mu := th.GetMu()
	mu.Lock()
	th.SetPendingHeaders(h)
	mu.Unlock()
	return nil
}

func whenBodyForProxyTransparency(th proxy.Transparency, doc *godog.DocString) error {
	if routerTest.srv == nil {
		return errors.New("router server is not started (call Given router is available first)")
	}
	if doc == nil {
		return errors.New("body doc string is missing")
	}

	body := []byte(strings.TrimSpace(doc.Content))

	mu := th.GetMu()
	mu.Lock()
	th.SetLastClientBody(append([]byte(nil), body...))
	var hdr http.Header
	if ph := th.GetPendingHeaders(); ph != nil {
		hdr = ph.Clone()
	}
	mu.Unlock()

	req, err := http.NewRequest(http.MethodPost, routerTest.srv.URL+config.ChatCompletionsUrl, bytes.NewReader(body))
	if err != nil {
		return err
	}
	for k, vals := range hdr {
		for _, v := range vals {
			req.Header.Add(k, v)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	routerTest.status = resp.StatusCode
	routerTest.body, err = io.ReadAll(resp.Body)
	return err
}

func jsonValuesEqualJSON(a, b []byte) (bool, error) {
	var av, bv any
	if err := json.Unmarshal(a, &av); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &bv); err != nil {
		return false, err
	}
	return reflect.DeepEqual(av, bv), nil
}

func thenUpstreamAReceivedSameJSONBody(th proxy.Transparency) error {
	mu := th.GetMu()
	mu.Lock()
	got := append([]byte(nil), th.GetLastUpstreamBody()...)
	want := append([]byte(nil), th.GetLastClientBody()...)
	mu.Unlock()
	if len(got) == 0 {
		return errors.New("upstream did not record a request body")
	}
	if len(want) == 0 {
		return errors.New("no client body was recorded for this scenario")
	}
	ok, err := jsonValuesEqualJSON(got, want)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return errors.New("upstream body JSON does not match the JSON body sent to the router")
}

func thenUpstreamAReceivedContentTypeApplicationJSON(th proxy.Transparency) error {
	mu := th.GetMu()
	mu.Lock()
	got := th.GetLastUpstreamCT()
	mu.Unlock()
	want := "application/json"
	if got != want {
		return errors.New(`upstream Content-Type: got ` + strconv.Quote(got) + `, want ` + strconv.Quote(want))
	}
	return nil
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

func cleanupProxyTransparencyScenario(th proxy.Transparency) func(context.Context, *godog.Scenario, error) (context.Context, error) {
	return func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		mu := th.GetMu()
		mu.Lock()
		if th.GetEchoSrv() != nil {
			th.GetEchoSrv().Close()
			th.SetEchoSrv(nil)
		}
		th.SetLastUpstreamBody(nil)
		th.SetLastUpstreamCT("")
		th.SetLastClientBody(nil)
		th.SetPendingHeaders(nil)
		mu.Unlock()
		resetRoundRobin()
		return ctx, nil
	}
}

func InitProxyTransparency(ctx *godog.ScenarioContext) {
	th := proxy.NewTransparency()

	ctx.Step(`^upstream A is configured to echo the received request for chat completions$`, func() error {
		return givenUpstreamEchoForChatCompletions(th)
	})
	ctx.Step(`^I send a POST request to "/v1/chat/completions" with headers:$`, func(table *godog.Table) error {
		return whenSendPostToChatCompletionsWithHeaders(th, table)
	})
	ctx.Step(`^body:$`, func(doc *godog.DocString) error {
		return whenBodyForProxyTransparency(th, doc)
	})
	ctx.Step(`^upstream A should have received the same JSON body$`, func() error {
		return thenUpstreamAReceivedSameJSONBody(th)
	})
	ctx.Step(`^upstream A should have received header "Content-Type" with value "application/json"$`, func() error {
		return thenUpstreamAReceivedContentTypeApplicationJSON(th)
	})

	ctx.After(cleanupProxyTransparencyScenario(th))
}
