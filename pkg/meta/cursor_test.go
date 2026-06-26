package meta

import (
	"testing"

	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
)

func TestFromProto(t *testing.T) {
	if FromProto(nil) != nil {
		t.Fatal("expected nil cursor meta")
	}

	limit := int32(20)
	got := FromProto(&sharedv1.CursorMeta{
		NextCursor: strPtr("next"),
		PrevCursor: strPtr("prev"),
		Limit:      &limit,
	})

	if got == nil {
		t.Fatal("expected cursor meta")
	}
	if got.NextCursor != "next" || got.PrevCursor != "prev" || got.Limit != 20 {
		t.Fatalf("unexpected cursor meta %#v", got)
	}
}

func strPtr(value string) *string {
	return &value
}
