package grpcclient

import (
	"errors"

	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapError(err error) error {
	if err == nil {
		return nil
	}

	var appErr *response.AppError
	if errors.As(err, &appErr) {
		return err
	}

	st, ok := status.FromError(err)
	if !ok {
		return response.Internal(err)
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return response.BadRequest("invalid_argument", st.Message(), err)
	case codes.Unauthenticated:
		return response.Unauthorized("unauthenticated", st.Message(), err)
	case codes.PermissionDenied:
		return response.Forbidden("permission_denied", st.Message(), err)
	case codes.NotFound:
		return response.NotFound("not_found", st.Message(), err)
	case codes.AlreadyExists:
		return response.Conflict("already_exists", st.Message(), err)
	case codes.FailedPrecondition:
		return response.FailedPrecondition("failed_precondition", st.Message(), err)
	default:
		return response.Internal(err)
	}
}
