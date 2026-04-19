package router

import (
	"bytes"
	"io"
	Dto "messhias/router-expirement/internal/DTO"
	"messhias/router-expirement/internal/balancer"
	"messhias/router-expirement/internal/proxy"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func registerChatRoutes(engine *gin.Engine, robin balancer.RoundRobin, hooks *proxy.Hooks) {
	engine.POST("/v1/chat/completions", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)

		c.Header("Content-Type", "application/json")

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if _, err := Dto.ParseAndValidateChatRequest(body); err != nil {
			c.Status(http.StatusBadRequest)
			_, _ = c.Writer.Write([]byte(`{"error":"invalid request"}`))
			return
		}

		target, err := robin.Next()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		upstreamURL := strings.TrimRight(target, "/") + "/v1/chat/completions"

		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, upstreamURL, bytes.NewReader(body))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if ct := c.GetHeader("Content-Type"); ct != "" {
			req.Header.Set("Content-Type", ct)
		} else {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.String(http.StatusBadGateway, err.Error())
			return
		}
		defer func() { _ = resp.Body.Close() }()

		upBody, err := io.ReadAll(resp.Body)
		if err != nil {
			c.String(http.StatusBadGateway, err.Error())
			return
		}

		if hooks != nil && hooks.OnUpstream2xx != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			hooks.OnUpstream2xx(strings.TrimRight(target, "/"))
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		c.Data(resp.StatusCode, contentType, upBody)
	})
}

func NewEngine(bal balancer.RoundRobin, hooks *proxy.Hooks) *gin.Engine {
	engine := gin.New()
	registerChatRoutes(engine, bal, hooks)
	return engine
}
