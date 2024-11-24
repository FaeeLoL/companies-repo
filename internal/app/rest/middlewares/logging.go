package middlewares

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type ctxKey string

const ctxKeyLogger ctxKey = "logger"

type LoggingMiddleware struct {
	logger *logrus.Logger
	next   http.Handler
}

func NewLoggingMiddleware(logger *logrus.Logger) func(r http.Handler) http.Handler {
	return func(r http.Handler) http.Handler {
		return &LoggingMiddleware{logger: logger, next: r}
	}
}

func (h *LoggingMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	startTime := time.Now()

	logger := h.logger.WithFields(
		map[string]any{
			"request_id":     generateOrExtractRequestID(r),
			"method":         r.Method,
			"uri":            r.URL.RequestURI(),
			"remote_addr":    r.RemoteAddr,
			"content_length": r.ContentLength,
		},
	)

	logger.Info("request started")

	r = r.WithContext(NewContextWithLogger(ctx, logger))
	wrw := wrapResponseWriterIfNeeded(rw)

	h.next.ServeHTTP(wrw, r)

	duration := time.Since(startTime)
	logger.WithFields(
		map[string]any{
			"duration_ms": duration.Milliseconds(),
			"status":      wrw.Status(),
			"bytes_sent":  wrw.BytesWritten(),
		}).Info(
		fmt.Sprintf("response completed in %.3fs", duration.Seconds()),
	)
}

func generateOrExtractRequestID(r *http.Request) string {
	if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
		return requestID
	}
	return generateSecureRequestID()
}

func generateSecureRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate random request ID")
	}
	return hex.EncodeToString(b)
}

// NewContextWithLogger creates a new context with logger.
func NewContextWithLogger(ctx context.Context, logger logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger, logger)
}

// GetLoggerFromContext extracts logger from the context.
func GetLoggerFromContext(ctx context.Context) logrus.FieldLogger {
	value := ctx.Value(ctxKeyLogger)
	if value == nil {
		return nil
	}
	return value.(logrus.FieldLogger)
}

func wrapResponseWriterIfNeeded(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	bytes       int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(buf []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	n, err := rw.ResponseWriter.Write(buf)
	rw.bytes += n
	return n, err
}

func (rw *responseWriter) Status() int {
	if !rw.wroteHeader {
		return http.StatusOK
	}
	return rw.status
}

func (rw *responseWriter) BytesWritten() int {
	return rw.bytes
}
