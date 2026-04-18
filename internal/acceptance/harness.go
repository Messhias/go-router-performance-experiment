package acceptance

import "errors"

type ChatAcceptanceHarness interface {
	EnsureTwoChatUpstreams() error
	Close() error
}

type chatAcceptanceHarness struct {
	// TODO: Add fields when you wire servers, e.g.:
	// TODO: upstreamA *httptest.Server
	// TODO: upstreamB *httptest.Server
}

func NewHarness() ChatAcceptanceHarness {
	return &chatAcceptanceHarness{}
}

func (h *chatAcceptanceHarness) EnsureTwoChatUpstreams() error {
	return errors.New("not implemented")
}

func (h *chatAcceptanceHarness) Close() error {
	// TODO: Close upstreamA/upstreamB when they exist; return joined errors if you want strict cleanup reporting.
	return nil
}
