package acceptance

import (
	"errors"
	"messhias/router-expirement/internal/upstreamfake"
	"net/http/httptest"
	"testing"
)

type ChatAcceptanceHarness interface {
	EnsureTwoChatUpstreams() error
	Close() error
}

type chatAcceptanceHarness struct {
	t         *testing.T
	upstreamA *httptest.Server
	upstreamB *httptest.Server
}

func NewHarness(t *testing.T) ChatAcceptanceHarness {
	t.Helper()
	return &chatAcceptanceHarness{
		t: t,
	}
}

func (h *chatAcceptanceHarness) EnsureTwoChatUpstreams() error {
	if h.upstreamA != nil {
		return errors.New("upstream server already set")
	}

	h.upstreamA = upstreamfake.NewChatCompletionServerMock(h.t, "upstream-a")
	h.upstreamB = upstreamfake.NewChatCompletionServerMock(h.t, "upstream-b")

	return nil
}

func (h *chatAcceptanceHarness) Close() error {
	// we do not call Close() here because in the completion serve we already call the close
	h.upstreamA = nil
	h.upstreamB = nil

	return nil
}
