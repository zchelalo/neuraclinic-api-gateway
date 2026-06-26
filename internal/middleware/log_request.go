package middleware

import (
	"net/http"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

type MetricsRecorder interface {
	Record(status int, duration time.Duration)
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func LogRequestMiddleware(metrics MetricsRecorder) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()
			recorder := &statusWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(recorder, r)

			if metrics != nil {
				metrics.Record(recorder.status, time.Since(startedAt))
			}

			logger := Logger(r.Context())
			if logger == nil {
				return
			}
			fields := []zap.Field{
				zap.Int("status", recorder.status),
				zap.Duration("duration", time.Since(startedAt)),
			}
			if recorder.status >= http.StatusBadRequest {
				logger.Warn("http request completed", fields...)
				return
			}
			logger.Info("http request completed", fields...)
		})
	}
}

type AtomicMetrics struct {
	Requests uint64
	Errors   uint64
}

func (m *AtomicMetrics) Record(status int, _ time.Duration) {
	atomic.AddUint64(&m.Requests, 1)
	if status >= http.StatusBadRequest {
		atomic.AddUint64(&m.Errors, 1)
	}
}
