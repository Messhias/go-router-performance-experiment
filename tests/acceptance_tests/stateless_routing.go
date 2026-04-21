package acceptance_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	Dto "messhias/router-expirement/internal/DTO"
	"messhias/router-expirement/internal/config"
	"net/http"
	"slices"
	"strings"

	"github.com/cucumber/godog"
)

func whenClientSendsPOSTWithHeaderAndBody(table *godog.Table) error {
	for _, row := range table.Rows[1:] {
		if len(row.Cells) < 2 {
			continue
		}

		clientId := strings.TrimSpace(row.Cells[0].Value)
		message := strings.TrimSpace(row.Cells[1].Value)

		body := Dto.ChatRequestDto{
			Model: "auto", Messages: []Dto.Message{
				{
					Role:    "user",
					Content: message,
				},
			},
		}

		bodyJson, err := json.Marshal(body)

		if err != nil {
			return err
		}

		req, err := http.NewRequest(
			http.MethodPost,
			routerTest.srv.URL+config.ChatCompletionsUrl,
			bytes.NewReader(bodyJson),
		)

		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Client-Id", clientId)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		b, readErr := io.ReadAll(resp.Body)
		closeErr := resp.Body.Close()

		if readErr != nil {
			return readErr
		}
		if closeErr != nil {
			return closeErr
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("status %d: %s", resp.StatusCode, string(b))
		}
	}

	return nil
}

func thenUpstreamHandlingOrderLast3ABABA() error {
	want := []string{"A", "B", "A"}
	roundRobinState.mu.Lock()
	defer roundRobinState.mu.Unlock()

	got := append([]string(nil), roundRobinState.handlingOrder...)

	got = got[len(got)-3:]

	if slices.Equal(got, want) == false {
		return fmt.Errorf("want %s, got %s", want, got)
	}

	return nil
}

func InitStatelessRouting(ctx *godog.ScenarioContext) {
	ctx.Step(`^following POST requests are sent to "/v1/chat/completions" in order with header "X-Client-Id" and JSON bodies built from:$`, whenClientSendsPOSTWithHeaderAndBody)
	ctx.Step(`^upstream handling order for the last 3 requests should be "A,B,A"$`, thenUpstreamHandlingOrderLast3ABABA)
}
