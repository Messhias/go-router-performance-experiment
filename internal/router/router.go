package router

import (
	"bytes"
	"io"
	Dto "messhias/router-expirement/internal/DTO"
	"messhias/router-expirement/internal/balancer"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func registerChatRoutes(engine *gin.Engine, robin balancer.RoundRobin) {
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

		upstreamUrl := strings.TrimRight(target, "/") + "/v1/chat/completions"

		req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, upstreamUrl, bytes.NewReader(body))

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

		contentType := resp.Header.Get("Content-Type")

		if contentType == "" {
			contentType = "application/octet-stream"
		}

		c.Data(resp.StatusCode, contentType, upBody)
	})
}

func NewEngine(balancer balancer.RoundRobin) *gin.Engine {
	engine := gin.New()

	registerChatRoutes(engine, balancer)

	return engine
}
