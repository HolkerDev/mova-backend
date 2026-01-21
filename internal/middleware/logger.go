package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
}))

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Read request body
		var requestBody []byte
		if c.Request.Body != nil {
			var err error
			requestBody, err = io.ReadAll(c.Request.Body)
			if err != nil {
				Logger.Error("failed to read request body", "error", err)
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Wrap response writer to capture response body
		rw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = rw

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		logAttrs := []any{
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", status),
			slog.Duration("duration", duration),
			slog.String("client_ip", c.ClientIP()),
		}

		// Log errors with request/response details
		if status >= 400 {
			logAttrs = append(logAttrs,
				slog.String("request_body", string(requestBody)),
				slog.String("response_body", rw.body.String()),
			)

			if status >= 500 {
				Logger.Error("request failed", logAttrs...)
			} else {
				Logger.Warn("request error", logAttrs...)
			}
		} else {
			Logger.Info("request completed", logAttrs...)
		}
	}
}
