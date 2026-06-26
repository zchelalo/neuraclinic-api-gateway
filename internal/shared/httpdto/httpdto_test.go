package httpdto

import (
	"testing"
	"time"

	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
	date "google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTimestampDateEnumAndStructHelpers(t *testing.T) {
	ts := timestamppb.New(time.Date(2026, 6, 25, 19, 0, 0, 0, time.UTC))
	gotTimestamp := Timestamp(ts)
	if gotTimestamp == nil || *gotTimestamp != "2026-06-25T19:00:00Z" {
		t.Fatalf("unexpected timestamp %#v", gotTimestamp)
	}

	gotDate := DateValue(&date.Date{Year: 2026, Month: 6, Day: 25})
	if gotDate == nil || gotDate.Year != 2026 || gotDate.Month != 6 || gotDate.Day != 25 {
		t.Fatalf("unexpected date %#v", gotDate)
	}

	if got := EnumString(sharedv1.RoleKey_ROLE_KEY_ADMIN); got != "ROLE_KEY_ADMIN" {
		t.Fatalf("unexpected enum %q", got)
	}
	if got := EnumString(sharedv1.RoleKey_ROLE_KEY_UNSPECIFIED); got != "" {
		t.Fatalf("expected empty zero enum, got %q", got)
	}

	value, err := structpb.NewStruct(map[string]any{"hello": "world"})
	if err != nil {
		t.Fatalf("new struct: %v", err)
	}
	gotStruct := Struct(value)
	if gotStruct["hello"] != "world" {
		t.Fatalf("unexpected struct map %#v", gotStruct)
	}
}
