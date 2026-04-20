package upstreamfake

import (
	"encoding/json"
	"messhias/router-expirement/internal/DTO"
	"messhias/router-expirement/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

// NewChatCompletionServerMock returns a httptest.Server that handles POST /v1/chat/completions.
// label is written into response JSON field "model" so tests can tell A from B.
// ATTENTION: DO NOT USE IT IN PRODUCTION
func NewChatCompletionServerMock(t *testing.T, label string) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc(config.ChatCompletionsUrl, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		resp := Dto.ChatCompletionResponseDto{
			ID:      "chatcmpl-test",
			Object:  "chat.completion",
			Created: 1,
			Model:   label,
			Choices: []Dto.ChoiceDto{
				{
					Index:        0,
					FinishReason: "stop",
					Message: Dto.MessageDto{
						Role:    "assistant",
						Content: "ok",
					},
				},
			},
		}

		payload, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}
