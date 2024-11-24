// Package middlewares contains HTTP middleware for request processing.
package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func NewErrorHandlerMiddleware(logger logrus.FieldLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					errMessage := fmt.Sprintf("%v", err)

					logger.WithFields(logrus.Fields{
						"path":   r.URL.Path,
						"method": r.Method,
						"error":  errMessage,
					}).Error("Unhandled error")

					rw.Header().Set("Content-Type", "application/json")
					rw.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(rw).Encode(map[string]string{
						"error": "Internal Server Error",
					})
				}
			}()
			next.ServeHTTP(rw, r)
		})
	}
}
