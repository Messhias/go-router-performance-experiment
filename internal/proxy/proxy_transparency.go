package proxy

import (
	"net/http"
	"net/http/httptest"
	"sync"
)

type Transparency interface {
	GetMu() *sync.Mutex

	GetPendingHeaders() http.Header
	SetPendingHeaders(h http.Header)

	GetEchoSrv() *httptest.Server
	SetEchoSrv(srv *httptest.Server)

	SetLastUpstreamBody(b []byte)
	GetLastUpstreamBody() []byte

	GetLastUpstreamCT() string
	SetLastUpstreamCT(ct string)

	GetLastClientBody() []byte
	SetLastClientBody(b []byte)
}

type transparency struct {
	mu sync.Mutex

	pendingHeaders http.Header
	echoSrv        *httptest.Server

	lastUpstreamBody []byte
	lastUpstreamCT   string
	lastClientBody   []byte
}

func (t *transparency) GetMu() *sync.Mutex {
	return &t.mu
}

func (t *transparency) GetPendingHeaders() http.Header {
	return t.pendingHeaders
}

func (t *transparency) SetPendingHeaders(h http.Header) {
	t.pendingHeaders = h
}

func (t *transparency) GetEchoSrv() *httptest.Server {
	return t.echoSrv
}

func (t *transparency) SetEchoSrv(srv *httptest.Server) {
	t.echoSrv = srv
}

func (t *transparency) SetLastUpstreamBody(b []byte) {
	t.lastUpstreamBody = b
}

func (t *transparency) GetLastUpstreamBody() []byte {
	return t.lastUpstreamBody
}

func (t *transparency) GetLastUpstreamCT() string {
	return t.lastUpstreamCT
}

func (t *transparency) SetLastUpstreamCT(ct string) {
	t.lastUpstreamCT = ct
}

func (t *transparency) GetLastClientBody() []byte {
	return t.lastClientBody
}

func (t *transparency) SetLastClientBody(b []byte) {
	t.lastClientBody = b
}

func NewTransparency() Transparency {
	return &transparency{}
}
