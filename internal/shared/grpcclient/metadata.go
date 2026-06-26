package grpcclient

import (
	"context"
	"net/http"
	"strings"

	"github.com/zchelalo/neuraclinic-api-gateway/internal/requestctx"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/language"
	"google.golang.org/grpc/metadata"
)

const (
	headerRequestID            = "x-request-id"
	headerTraceID              = "x-trace-id"
	headerAcceptLanguage       = "accept-language"
	headerUserID               = "x-user-id"
	headerPsychologistID       = "x-psychologist-id"
	headerAdminID              = "x-admin-id"
	headerInternalServiceToken = "x-internal-service-token"
)

type CallOptions struct {
	IncludeAuth          bool
	InternalServiceToken string
}

func OutgoingContext(ctx context.Context, r *http.Request, opts CallOptions) context.Context {
	md := metadata.New(nil)
	if requestID := requestctx.RequestID(ctx); requestID != "" {
		md.Set(headerRequestID, requestID)
	}
	if traceID := requestctx.TraceID(ctx); traceID != "" {
		md.Set(headerTraceID, traceID)
	}
	md.Set(headerAcceptLanguage, language.ResolveRequest(r))
	if opts.IncludeAuth {
		if auth, ok := requestctx.Auth(ctx); ok {
			if auth.UserID != "" {
				md.Set(headerUserID, auth.UserID)
			}
			if auth.PsychologistID != "" {
				md.Set(headerPsychologistID, auth.PsychologistID)
			}
			if auth.AdminID != "" {
				md.Set(headerAdminID, auth.AdminID)
			}
		}
	}
	if strings.TrimSpace(opts.InternalServiceToken) != "" {
		md.Set(headerInternalServiceToken, strings.TrimSpace(opts.InternalServiceToken))
	}
	return metadata.NewOutgoingContext(ctx, md)
}
