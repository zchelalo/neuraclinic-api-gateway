package middleware

import (
	"context"

	"github.com/zchelalo/neuraclinic-api-gateway/internal/requestctx"
	"go.uber.org/zap"
)

type AuthMode = requestctx.AuthMode

const (
	AuthModeWeb    = requestctx.AuthModeWeb
	AuthModeMobile = requestctx.AuthModeMobile
)

type AuthContext = requestctx.AuthContext

func WithRequestID(ctx context.Context, value string) context.Context {
	return requestctx.WithRequestID(ctx, value)
}

func RequestID(ctx context.Context) string {
	return requestctx.RequestID(ctx)
}

func WithTraceID(ctx context.Context, value string) context.Context {
	return requestctx.WithTraceID(ctx, value)
}

func TraceID(ctx context.Context) string {
	return requestctx.TraceID(ctx)
}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return requestctx.WithLogger(ctx, logger)
}

func Logger(ctx context.Context) *zap.Logger {
	return requestctx.Logger(ctx)
}

func WithAuth(ctx context.Context, auth AuthContext) context.Context {
	return requestctx.WithAuth(ctx, auth)
}

func Auth(ctx context.Context) (AuthContext, bool) {
	return requestctx.Auth(ctx)
}
