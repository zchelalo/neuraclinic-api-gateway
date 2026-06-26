package httpdto

import (
	"time"

	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
	date "google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Date struct {
	Year  int32 `json:"year,omitempty"`
	Month int32 `json:"month,omitempty"`
	Day   int32 `json:"day,omitempty"`
}

type OperationResponse struct {
	Message  string            `json:"message,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type enumStringer interface {
	~int32
	String() string
}

func Timestamp(value *timestamppb.Timestamp) *string {
	if value == nil {
		return nil
	}
	formatted := value.AsTime().UTC().Format(time.RFC3339Nano)
	return &formatted
}

func DateValue(value *date.Date) *Date {
	if value == nil {
		return nil
	}
	return &Date{
		Year:  value.GetYear(),
		Month: value.GetMonth(),
		Day:   value.GetDay(),
	}
}

func Operation(value *sharedv1.OperationResponse) *OperationResponse {
	if value == nil {
		return nil
	}
	return &OperationResponse{
		Message:  value.GetMessage(),
		Metadata: value.GetMetadata(),
	}
}

func Struct(value *structpb.Struct) map[string]any {
	if value == nil {
		return nil
	}
	return value.AsMap()
}

func EnumString[T enumStringer](value T) string {
	if value == 0 {
		return ""
	}
	return value.String()
}
