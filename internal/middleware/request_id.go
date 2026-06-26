package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

const (
	headerRequestID = "X-Request-Id"
	headerTraceID   = "X-Trace-Id"
)

func RequestIDMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(headerRequestID)
			if requestID == "" {
				requestID = uuid.NewString()
			}
			traceID := r.Header.Get(headerTraceID)
			if traceID == "" {
				traceID = uuid.NewString()
			}

			ctx := WithRequestID(r.Context(), requestID)
			ctx = WithTraceID(ctx, traceID)
			r = r.WithContext(ctx)

			w.Header().Set(headerRequestID, requestID)
			w.Header().Set(headerTraceID, traceID)

			next.ServeHTTP(w, r)
		})
	}
}
