package requestctx

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	traceIDKey   contextKey = "trace_id"
	loggerKey    contextKey = "logger"
	authKey      contextKey = "auth"
)

type AuthMode string

const (
	AuthModeWeb    AuthMode = "web"
	AuthModeMobile AuthMode = "mobile"
)

type AuthContext struct {
	Token           string
	Mode            AuthMode
	UserID          string
	RoleKey         string
	PsychologistID  string
	AdminID         string
	PermissionsKeys []string
}

func WithRequestID(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, requestIDKey, value)
}

func RequestID(ctx context.Context) string {
	value, _ := ctx.Value(requestIDKey).(string)
	return value
}

func WithTraceID(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, traceIDKey, value)
}

func TraceID(ctx context.Context) string {
	value, _ := ctx.Value(traceIDKey).(string)
	return value
}

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func Logger(ctx context.Context) *zap.Logger {
	logger, _ := ctx.Value(loggerKey).(*zap.Logger)
	return logger
}

func WithAuth(ctx context.Context, auth AuthContext) context.Context {
	return context.WithValue(ctx, authKey, auth)
}

func Auth(ctx context.Context) (AuthContext, bool) {
	value, ok := ctx.Value(authKey).(AuthContext)
	return value, ok
}
