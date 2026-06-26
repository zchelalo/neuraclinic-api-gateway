package parse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	locationv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/location/v1"
	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func JSONBody[T any](r *http.Request, dst *T) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return response.BadRequest("invalid_json", "invalid json body", err)
	}
	return nil
}

func QueryInt(r *http.Request, key string, fallback int32) (int32, error) {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, response.BadRequest("invalid_query", fmt.Sprintf("query param %s must be an integer", key), err)
	}
	return int32(value), nil
}

func QueryString(r *http.Request, key string) string {
	return strings.TrimSpace(r.URL.Query().Get(key))
}

func QueryOptionalBool(r *http.Request, key string) (*bool, error) {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return nil, nil
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return nil, response.BadRequest("invalid_query", fmt.Sprintf("query param %s must be a boolean", key), err)
	}
	return &value, nil
}

func OptionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func Timestamp(value string) (*timestamppb.Timestamp, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, response.BadRequest("invalid_timestamp", "timestamp must be RFC3339", err)
	}
	return timestamppb.New(parsed), nil
}

func Date(value string) (*date.Date, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, response.BadRequest("invalid_date", "date must use YYYY-MM-DD", err)
	}
	return &date.Date{
		Year:  int32(parsed.Year()),
		Month: int32(parsed.Month()),
		Day:   int32(parsed.Day()),
	}, nil
}

func CursorPagination(r *http.Request, defaultLimit int32) (*sharedv1.CursorPagination, error) {
	limit, err := QueryInt(r, "limit", defaultLimit)
	if err != nil {
		return nil, err
	}
	return &sharedv1.CursorPagination{
		AfterCursor:  OptionalString(QueryString(r, "after_cursor")),
		BeforeCursor: OptionalString(QueryString(r, "before_cursor")),
		Limit:        limit,
	}, nil
}

func DateRange(startRaw, endRaw string) (*sharedv1.DateRange, error) {
	start, err := Timestamp(startRaw)
	if err != nil {
		return nil, err
	}
	end, err := Timestamp(endRaw)
	if err != nil {
		return nil, err
	}
	if start == nil && end == nil {
		return nil, nil
	}
	if start == nil || end == nil {
		return nil, response.BadRequest("invalid_date_range", "date range requires start_date and end_date", nil)
	}
	return &sharedv1.DateRange{
		StartDate: start,
		EndDate:   end,
	}, nil
}

func RoleKey(value string) (sharedv1.RoleKey, error) {
	switch normalizeEnum(value) {
	case "ROLE_KEY_ADMIN", "ADMIN":
		return sharedv1.RoleKey_ROLE_KEY_ADMIN, nil
	case "ROLE_KEY_PSYCHOLOGIST", "PSYCHOLOGIST":
		return sharedv1.RoleKey_ROLE_KEY_PSYCHOLOGIST, nil
	default:
		return sharedv1.RoleKey_ROLE_KEY_UNSPECIFIED, response.BadRequest("invalid_role_key", "invalid role_key", nil)
	}
}

func Sex(value string) (recordv1.Sex, error) {
	switch normalizeEnum(value) {
	case "SEX_MALE", "MALE":
		return recordv1.Sex_SEX_MALE, nil
	case "SEX_FEMALE", "FEMALE":
		return recordv1.Sex_SEX_FEMALE, nil
	case "SEX_OTHER", "OTHER":
		return recordv1.Sex_SEX_OTHER, nil
	case "SEX_PREFER_NOT_TO_SAY", "PREFER_NOT_TO_SAY":
		return recordv1.Sex_SEX_PREFER_NOT_TO_SAY, nil
	default:
		return recordv1.Sex_SEX_UNSPECIFIED, response.BadRequest("invalid_sex", "invalid sex", nil)
	}
}

func MaritalStatus(value string) (recordv1.MaritalStatus, error) {
	switch normalizeEnum(value) {
	case "MARITAL_STATUS_SINGLE", "SINGLE":
		return recordv1.MaritalStatus_MARITAL_STATUS_SINGLE, nil
	case "MARITAL_STATUS_MARRIED", "MARRIED":
		return recordv1.MaritalStatus_MARITAL_STATUS_MARRIED, nil
	case "MARITAL_STATUS_DIVORCED", "DIVORCED":
		return recordv1.MaritalStatus_MARITAL_STATUS_DIVORCED, nil
	case "MARITAL_STATUS_WIDOWED", "WIDOWED":
		return recordv1.MaritalStatus_MARITAL_STATUS_WIDOWED, nil
	case "MARITAL_STATUS_SEPARATED", "SEPARATED":
		return recordv1.MaritalStatus_MARITAL_STATUS_SEPARATED, nil
	case "MARITAL_STATUS_COHABITING", "COHABITING":
		return recordv1.MaritalStatus_MARITAL_STATUS_COHABITING, nil
	default:
		return recordv1.MaritalStatus_MARITAL_STATUS_UNSPECIFIED, response.BadRequest("invalid_marital_status", "invalid marital_status", nil)
	}
}

func AppointmentStatus(value string) (sharedv1.AppointmentStatus, error) {
	switch normalizeEnum(value) {
	case "APPOINTMENT_STATUS_SCHEDULED", "SCHEDULED":
		return sharedv1.AppointmentStatus_APPOINTMENT_STATUS_SCHEDULED, nil
	case "APPOINTMENT_STATUS_COMPLETED", "COMPLETED":
		return sharedv1.AppointmentStatus_APPOINTMENT_STATUS_COMPLETED, nil
	case "APPOINTMENT_STATUS_RESCHEDULED", "RESCHEDULED":
		return sharedv1.AppointmentStatus_APPOINTMENT_STATUS_RESCHEDULED, nil
	case "APPOINTMENT_STATUS_CANCELLED", "CANCELLED":
		return sharedv1.AppointmentStatus_APPOINTMENT_STATUS_CANCELLED, nil
	case "APPOINTMENT_STATUS_NO_SHOW", "NO_SHOW":
		return sharedv1.AppointmentStatus_APPOINTMENT_STATUS_NO_SHOW, nil
	default:
		return sharedv1.AppointmentStatus_APPOINTMENT_STATUS_UNSPECIFIED, response.BadRequest("invalid_status", "invalid appointment status", nil)
	}
}

func FileStatus(value string) (sharedv1.FileStatus, error) {
	switch normalizeEnum(value) {
	case "FILE_STATUS_UPLOADING", "UPLOADING":
		return sharedv1.FileStatus_FILE_STATUS_UPLOADING, nil
	case "FILE_STATUS_AVAILABLE", "AVAILABLE":
		return sharedv1.FileStatus_FILE_STATUS_AVAILABLE, nil
	case "FILE_STATUS_DELETED", "DELETED":
		return sharedv1.FileStatus_FILE_STATUS_DELETED, nil
	case "FILE_STATUS_ERROR", "ERROR":
		return sharedv1.FileStatus_FILE_STATUS_ERROR, nil
	default:
		return sharedv1.FileStatus_FILE_STATUS_UNSPECIFIED, response.BadRequest("invalid_status", "invalid file status", nil)
	}
}

func AdminAreaType(value string) (*locationv1.AdminAreaType, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	var parsed locationv1.AdminAreaType
	switch normalizeEnum(value) {
	case "ADMIN_AREA_TYPE_STATE", "STATE":
		parsed = locationv1.AdminAreaType_ADMIN_AREA_TYPE_STATE
	case "ADMIN_AREA_TYPE_MUNICIPALITY", "MUNICIPALITY":
		parsed = locationv1.AdminAreaType_ADMIN_AREA_TYPE_MUNICIPALITY
	default:
		return nil, response.BadRequest("invalid_admin_area_type", "invalid admin_area_type", nil)
	}
	return &parsed, nil
}

func normalizeEnum(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "-", "_")
	value = strings.ReplaceAll(value, " ", "_")
	return strings.ToUpper(value)
}
