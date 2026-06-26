package grpcclient

import (
	"net/http"
	"testing"

	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapError(t *testing.T) {
	tests := []struct {
		name   string
		code   codes.Code
		status int
	}{
		{name: "invalid_argument", code: codes.InvalidArgument, status: http.StatusBadRequest},
		{name: "unauthenticated", code: codes.Unauthenticated, status: http.StatusUnauthorized},
		{name: "permission_denied", code: codes.PermissionDenied, status: http.StatusForbidden},
		{name: "not_found", code: codes.NotFound, status: http.StatusNotFound},
		{name: "already_exists", code: codes.AlreadyExists, status: http.StatusConflict},
		{name: "failed_precondition", code: codes.FailedPrecondition, status: http.StatusPreconditionFailed},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := MapError(status.Error(tc.code, tc.name))
			appErr, ok := err.(*response.AppError)
			if !ok {
				t.Fatalf("expected AppError, got %T", err)
			}
			if appErr.Status != tc.status {
				t.Fatalf("expected status %d, got %d", tc.status, appErr.Status)
			}
		})
	}
}
