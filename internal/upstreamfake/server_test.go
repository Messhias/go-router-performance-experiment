package upstreamfake

import (
	"bytes"
	"encoding/json"
	"io"
	Dto "messhias/router-expirement/internal/DTO"
	"messhias/router-expirement/internal/config"
	"net/http"
	"testing"
)

func TestUpstreams_ShouldPass(t *testing.T) {

	label := "upstream-a"
	runUpStreamTest(t, label)

	label = "upstream-b"
	runUpStreamTest(t, label)
}

func TestWrongMethodInWrongUrl_ShouldFail(t *testing.T) {
	srv := NewChatCompletionServerMock(t, "A")

	req, err := http.NewRequest(http.MethodGet, srv.URL+config.ChatCompletionsUrl, nil)

	res, err := srv.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(res.Body)

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("status=%d want %d", res.StatusCode, http.StatusMethodNotAllowed)
	}
}

func runUpStreamTest(t *testing.T, label string) {
	srv := NewChatCompletionServerMock(t, label)
	body, err := json.Marshal(Dto.ChatRequestDto{
		Model:    "auto",
		Messages: []Dto.Message{{Role: "user", Content: "Hello"}},
	})

	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, srv.URL+config.ChatCompletionsUrl, bytes.NewReader(body))

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := srv.Client().Do(req)

	if err != nil {
		t.Fatal(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		t.Fatalf("status=%d body=%s", res.StatusCode, b)
	}

	var got Dto.ChatResponseDto

	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	if got.Object != "chat.completion" {
		t.Fatalf("object=%q", got.Object)
	}

	if got.Model != label {
		t.Fatalf("model=%q want %s", got.Model, label)
	}
}
