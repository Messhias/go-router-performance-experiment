package router

import (
	"bytes"
	"context"
	"errors"
	"io"
	Dto "messhias/router-expirement/internal/DTO"
	"messhias/router-expirement/internal/balancer"
	"messhias/router-expirement/internal/config"
	"messhias/router-expirement/internal/proxy"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var upstreamRequestTimeout = 2 * time.Second

func registerChatRoutes(engine *gin.Engine, robin balancer.RoundRobin, hooks *proxy.Hooks) {
	engine.POST(config.ChatCompletionsUrl, func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		if _, err := Dto.ParseAndValidateChatRequest(body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		target, err := robin.Next()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		upstreamURL := strings.TrimRight(target, "/") + config.ChatCompletionsUrl

		ctx, cancel := context.WithTimeout(c.Request.Context(), upstreamRequestTimeout)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, upstreamURL, bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		if ct := c.GetHeader("Content-Type"); ct != "" {
			req.Header.Set("Content-Type", ct)
		} else {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {

			if errors.Is(err, context.DeadlineExceeded) || isTimeoutNetError(err) {
				c.JSON(http.StatusBadGateway, gin.H{"error": "upstream timeout"})
				return
			}

			c.JSON(http.StatusBadGateway, gin.H{"error": "upstream unavailable"})
			return
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode >= 500 {
			c.JSON(http.StatusBadGateway, gin.H{"error": "upstream unavailable"})
			return
		}

		upBody, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "upstream unavailable"})
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

func isTimeoutNetError(err error) bool {
	var ne net.Error
	return errors.As(err, &ne) && ne.Timeout()
}

func NewEngine(bal balancer.RoundRobin, hooks *proxy.Hooks) *gin.Engine {
	engine := gin.New()
	registerChatRoutes(engine, bal, hooks)
	return engine
}
