package gin

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/ZyoGo/default-ddd-http/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type bodyDumpResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

type Skipper func(c *gin.Context) bool

func ZerologLoggerWithSkipper(log zerolog.Logger, skipper Skipper) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Get request body
		reqBody := []byte{}
		contentType := c.GetHeader("Content-Type")
		if c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
		}
		mapData := make(map[string]interface{})

		if len(reqBody) > 0 && contentType == "application/json" {
			if err := json.Unmarshal(reqBody, &mapData); err != nil {
				c.Error(err)
			}
		}
		// Reset request body
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		// Masking credentials field
		doc := &Document{}
		bodyMasked := doc.throughMap(mapData)

		// Get response body
		resBody := new(bytes.Buffer)
		writer := &bodyDumpResponseWriter{body: resBody, ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		// Skip log
		if skipper != nil && skipper(c) {
			return
		}

		resBodyMap := make(map[string]interface{})
		if len(resBody.Bytes()) > 0 {
			if err := json.Unmarshal(resBody.Bytes(), &resBodyMap); err != nil {
				c.Error(err)
			}
		}
		resBodyMasked := doc.throughMap(resBodyMap)

		// Logging fields
		latency := time.Since(start)
		status := c.Writer.Status()
		req := c.Request

		logEvent := logger.Get().With().
			Int("status", status).
			Str("latency", latency.String()).
			Str("method", req.Method).
			Str("uri", req.RequestURI).
			Str("host", req.Host).
			Str("remote_ip", c.ClientIP()).
			Interface("headers", req.Header).
			Interface("request_body", bodyMasked).
			Interface("response_body", resBodyMasked).
			Logger()

		// Check request ID in headers
		id := c.GetHeader("X-Request-ID")
		if id != "" {
			logEvent = logEvent.With().Str("id", id).Logger()
		}

		switch {
		case status >= 500:
			logEvent.Error().Msg("Server Error")
		case status >= 400:
			logEvent.Warn().Msg("Client Error")
		case status >= 300:
			logEvent.Info().Msg("Redirection")
		default:
			logEvent.Info().Msg("Success")
		}
	}
}
