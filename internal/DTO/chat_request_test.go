package Dto

import (
	"testing"
)

func TestChatRequestDto_ShouldPass(t *testing.T) {
	validMinimal(t)
	invalidJson(t)
	validate(t)
}

func validate(t *testing.T) {
	tests := []struct {
		name string
		body []byte
	}{
		{
			name: "empty model and empty message fields",
			body: []byte(`{"model":"","messages":[{"role":"","content":""}]}`),
		},
		{
			name: "empty model empty content",
			body: []byte(`{"model":"x","messages":[{"role":"x","content":""}]}`),
		},
		{
			name: "empty model empty role",
			body: []byte(`{"model":"x","messages":[{"role":"","content":"p"}]}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := ParseChatRequest(tt.body)
			if err != nil {
				t.Fatal(err)
			}

			err = parsed.Validate()

			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func invalidJson(t *testing.T) bool {
	return t.Run("invalid json", func(t *testing.T) {
		_, err := ParseAndValidateChatRequest([]byte(`{`))
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func validMinimal(t *testing.T) bool {
	return t.Run("valid minimal", func(t *testing.T) {
		body := []byte(`{"model":"auto","messages":[{"role":"user","content":"Hi"}]}`)
		req, err := ParseAndValidateChatRequest(body)
		if err != nil {
			t.Fatal(err)
		}
		if req.Model != "auto" {
			t.Fatalf("model=%q", req.Model)
		}
	})
}
