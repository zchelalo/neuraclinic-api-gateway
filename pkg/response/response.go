package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/zchelalo/neuraclinic-api-gateway/internal/requestctx"
	"google.golang.org/protobuf/proto"
)

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type Envelope struct {
	Data      any           `json:"data,omitempty"`
	Meta      any           `json:"meta,omitempty"`
	Error     *ErrorPayload `json:"error,omitempty"`
	RequestID string        `json:"request_id"`
}

type AppError struct {
	Status  int
	Code    string
	Message string
	Details any
	Err     error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func NewError(status int, code, message string, err error) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func BadRequest(code, message string, err error) *AppError {
	return NewError(http.StatusBadRequest, code, message, err)
}

func Unauthorized(code, message string, err error) *AppError {
	return NewError(http.StatusUnauthorized, code, message, err)
}

func Forbidden(code, message string, err error) *AppError {
	return NewError(http.StatusForbidden, code, message, err)
}

func NotFound(code, message string, err error) *AppError {
	return NewError(http.StatusNotFound, code, message, err)
}

func Conflict(code, message string, err error) *AppError {
	return NewError(http.StatusConflict, code, message, err)
}

func FailedPrecondition(code, message string, err error) *AppError {
	return NewError(http.StatusPreconditionFailed, code, message, err)
}

func Internal(err error) *AppError {
	return NewError(http.StatusInternalServerError, "internal_error", "internal server error", err)
}

func WithDetails(err *AppError, details any) *AppError {
	if err == nil {
		return nil
	}
	err.Details = details
	return err
}

func Write(w http.ResponseWriter, r *http.Request, status int, data any, meta any) {
	mustNotContainProto(data)
	mustNotContainProto(meta)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Envelope{
		Data:      data,
		Meta:      meta,
		RequestID: requestctx.RequestID(r.Context()),
	})
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	appErr := Internal(err)
	var candidate *AppError
	if errors.As(err, &candidate) && candidate != nil {
		appErr = candidate
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Status)
	mustNotContainProto(appErr.Details)
	_ = json.NewEncoder(w).Encode(Envelope{
		Error: &ErrorPayload{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		},
		RequestID: requestctx.RequestID(r.Context()),
	})
}

func mustNotContainProto(value any) {
	visited := map[uintptr]struct{}{}
	if containsProto(reflect.ValueOf(value), visited) {
		panic("response.Write does not support protobuf payloads")
	}
}

func containsProto(value reflect.Value, visited map[uintptr]struct{}) bool {
	if !value.IsValid() {
		return false
	}

	if value.CanInterface() {
		if _, ok := value.Interface().(proto.Message); ok {
			return true
		}
	}

	switch value.Kind() {
	case reflect.Interface:
		if value.IsNil() {
			return false
		}
		return containsProto(value.Elem(), visited)
	case reflect.Pointer:
		if value.IsNil() {
			return false
		}
		ptr := value.Pointer()
		if ptr != 0 {
			if _, ok := visited[ptr]; ok {
				return false
			}
			visited[ptr] = struct{}{}
		}
		return containsProto(value.Elem(), visited)
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			if containsProto(value.Index(i), visited) {
				return true
			}
		}
		return false
	case reflect.Map:
		iter := value.MapRange()
		for iter.Next() {
			if containsProto(iter.Value(), visited) {
				return true
			}
		}
		return false
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			field := value.Type().Field(i)
			if !field.IsExported() {
				continue
			}
			if containsProto(value.Field(i), visited) {
				return true
			}
		}
		return false
	default:
		return false
	}
}
