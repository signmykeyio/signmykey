package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

// RequestLoggerKey is the key that holds the logger instance in a request context
var RequestLoggerKey = contextKey("logger")

// StatusRecorder is an object allowing http response status code to be logger
type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

// WriteHeader set StatusRecorder Status field with http response status code
func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

// Logger is a middleware for http logging of server/api requests
func Logger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()

			recorder := &StatusRecorder{
				ResponseWriter: w,
				Status:         200,
			}

			reqID := middleware.GetReqID(r.Context())

			logger.WithFields(logrus.Fields{
				"ctx":    "http",
				"host":   r.Host,
				"path":   r.URL.Path,
				"ip":     r.RemoteAddr,
				"proto":  r.Proto,
				"method": r.Method,
				"req_id": reqID,
			}).Info("HTTP Request")

			defer func() {
				logger.WithFields(logrus.Fields{
					"ctx":      "http",
					"host":     r.Host,
					"path":     r.URL.Path,
					"ip":       r.RemoteAddr,
					"proto":    r.Proto,
					"method":   r.Method,
					"duration": time.Since(t1).String(),
					"status":   recorder.Status,
					"req_id":   reqID,
				}).Info("HTTP Response")
			}()

			// Embed logger in context
			ctx := context.WithValue(r.Context(), RequestLoggerKey, logger)

			next.ServeHTTP(recorder, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
