package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func LoggerMiddleware(baseLogger *zap.Logger, serviceName string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := baseLogger.With(
				zap.String("service", serviceName),
				zap.String("request_id", RequestID(r.Context())),
				zap.String("trace_id", TraceID(r.Context())),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
			)
			next.ServeHTTP(w, r.WithContext(WithLogger(r.Context(), logger)))
		})
	}
}
