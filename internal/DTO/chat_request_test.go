package Dto

import "testing"

func TestParseAndValidateChatRequest_ShouldPass(t *testing.T) {
	t.Run("valid minimal", func(t *testing.T) {
		body := []byte(`{"model":"auto","messages":[{"role":"user","content":"Hi"}]}`)
		req, err := ParseAndValidateChatRequest(body)
		if err != nil {
			t.Fatal(err)
		}
		if req.Model != "auto" {
			t.Fatalf("model=%q", req.Model)
		}
	})
	t.Run("invalid json", func(t *testing.T) {
		_, err := ParseAndValidateChatRequest([]byte(`{`))
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
