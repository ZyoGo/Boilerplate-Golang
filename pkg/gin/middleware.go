package gin

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/ZyoGo/default-ddd-http/pkg/logger"
	"github.com/gin-gonic/gin"
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

func ZerologLoggerWithSkipper(skipper Skipper) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Capture request body
		reqBody := captureRequestBody(c)
		maskedReqBody := maskSensitiveFields(reqBody)

		// Capture response body
		resBody := new(bytes.Buffer)
		writer := &bodyDumpResponseWriter{body: resBody, ResponseWriter: c.Writer}
		c.Writer = writer

		// Process the request
		c.Next()

		// Skip logging if the skipper returns true
		if skipper != nil && skipper(c) {
			return
		}

		// Capture and mask response body
		resBodyMap := parseJSON(resBody.Bytes())
		maskedResBody := maskSensitiveFields(resBodyMap)

		// Log the request and response
		logRequestResponse(c, start, maskedReqBody, maskedResBody)
	}
}

func captureRequestBody(c *gin.Context) map[string]interface{} {
	var body []byte
	if c.Request.Body != nil {
		body, _ = io.ReadAll(c.Request.Body)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	contentType := c.GetHeader("Content-Type")
	if contentType == "application/json" && len(body) > 0 {
		return parseJSON(body)
	}
	return nil
}

func parseJSON(data []byte) map[string]interface{} {
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil
	}
	return result
}

func maskSensitiveFields(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		return nil
	}

	// Assuming Document is a struct that has a method throughMap
	doc := &Document{}
	return doc.ProcessMap(data)
}

func logRequestResponse(c *gin.Context, start time.Time, reqBody, resBody map[string]interface{}) {
	status := c.Writer.Status()
	latency := time.Since(start)
	req := c.Request

	logEvent := logger.Get().With().
		Int("status", status).
		Str("latency", latency.String()).
		Str("method", req.Method).
		Str("uri", req.RequestURI).
		Str("host", req.Host).
		Str("remote_ip", c.ClientIP()).
		Interface("headers", req.Header).
		Interface("request_body", reqBody).
		Interface("response_body", resBody).
		Logger()

	if id := c.GetHeader("X-Request-ID"); id != "" {
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
