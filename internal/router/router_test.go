package router

import (
	"bytes"
	"encoding/json"
	Dto "messhias/router-expirement/internal/DTO"
	"messhias/router-expirement/internal/balancer"
	"messhias/router-expirement/internal/upstreamfake"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func TestPOSTChatCompletions_proxiesToUpstream_openAICompatibleJSON_ShouldPass(t *testing.T) {
	urlA := upstreamfake.NewChatCompletionServerMock(t, "upstream-a")
	urlB := upstreamfake.NewChatCompletionServerMock(t, "upstream-b")

	bal, err := balancer.NewBalancer([]string{urlA.URL, urlB.URL})

	if err != nil {
		t.Fatalf("error creating balancer: %v", err)
	}

	engine := NewEngine(bal)
	body := []byte(`{
		"model": "auto",
		"messages": [{"role": "user", "content": "Hello"}]
	}`)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", w.Code, http.StatusOK)
	}

	var parsed Dto.ChatCompletionResponseDto

	if err := json.Unmarshal(w.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("error unmarshalling response body: %v", err)
	}

	if parsed.Object != "chat.completion" {
		t.Fatalf("POST /v1/chat/completions returned %s", parsed.Object)
	}

	if len(parsed.Choices) == 0 {
		t.Fatalf("POST /v1/chat/completions returned no choices")
	}

	if parsed.Choices[0].Message.Role != "assistant" {
		t.Fatalf("POST /v1/chat/completions returned wrong role")
	}

	if parsed.Model != "upstream-a" {
		t.Fatalf("model = %q, want upstream-a (first RR target)", parsed.Model)
	}
}
