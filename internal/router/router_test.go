package router

import (
	"bytes"
	"encoding/json"
	"io"
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

	engine := NewEngine(bal, nil)
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

func TestPOSTChatCompletions_upstreamReceivesPostMethod_ShouldPass(t *testing.T) {
	var gotMethod string

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		_, _ = io.Copy(io.Discard, r.Body)

		resp := Dto.ChatCompletionResponseDto{
			ID:      "chatcmpl-test",
			Object:  "chat.completion",
			Created: 1,
			Model:   "upstream",
			Choices: []Dto.ChoiceDto{{
				Index:        0,
				FinishReason: "stop",
				Message:      Dto.MessageDto{Role: "assistant", Content: "ok"},
			}},
		}
		payload, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
	}))
	t.Cleanup(upstream.Close)

	bal, err := balancer.NewBalancer([]string{upstream.URL})
	if err != nil {
		t.Fatalf("balancer: %v", err)
	}

	body := []byte(`{"model":"auto","messages":[{"role":"user","content":"Hello"}]}`)
	engine := NewEngine(bal, nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("upstream method = %q, want POST", gotMethod)
	}
}

func TestPOSTChatCompletions_upstreamReceivesSameBodyAndContentType_ShouldPass(t *testing.T) {
	var gotBody []byte
	var gotCT string

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		gotBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read body: %v", err)
		}
		gotCT = r.Header.Get("Content-Type")

		resp := Dto.ChatCompletionResponseDto{
			ID:      "chatcmpl-test",
			Object:  "chat.completion",
			Created: 1,
			Model:   "upstream",
			Choices: []Dto.ChoiceDto{{
				Index:        0,
				FinishReason: "stop",
				Message:      Dto.MessageDto{Role: "assistant", Content: "ok"},
			}},
		}
		payload, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
	}))
	t.Cleanup(upstream.Close)

	bal, err := balancer.NewBalancer([]string{upstream.URL})
	if err != nil {
		t.Fatalf("balancer: %v", err)
	}

	wantBody := []byte(`{"model":"auto","messages":[{"role":"user","content":"Hello"}]}`)
	engine := NewEngine(bal, nil)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", bytes.NewReader(wantBody))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	engine.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if !bytes.Equal(gotBody, wantBody) {
		t.Fatalf("upstream body = %q, want %q", gotBody, wantBody)
	}
	if gotCT != "application/json; charset=utf-8" {
		t.Fatalf("upstream Content-Type = %q, want application/json; charset=utf-8", gotCT)
	}
}
