package v1

import (
	"testing"
	"time"

	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
	userv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFromProtoUserFindByIDResponsePreservesOneofAndOptionalFields(t *testing.T) {
	now := timestamppb.New(time.Date(2026, 6, 25, 13, 0, 0, 0, time.UTC))
	resp := fromProtoUserFindByIDResponse(&userv1.UserServiceFindByIdResponse{
		User: &userv1.User{
			Id:        "user-1",
			Email:     "admin@example.com",
			RoleKey:   sharedv1.RoleKey_ROLE_KEY_ADMIN,
			UpdatedAt: now,
		},
		Profile: &userv1.UserServiceFindByIdResponse_Admin{
			Admin: &userv1.AdminProfile{Id: "admin-1"},
		},
	})

	if resp == nil || resp.User == nil {
		t.Fatal("expected mapped user response")
	}
	if resp.Admin == nil {
		t.Fatal("expected admin profile")
	}
	if resp.Psychologist != nil {
		t.Fatal("expected psychologist profile to stay nil")
	}
	if resp.User.RoleKey != "ROLE_KEY_ADMIN" {
		t.Fatalf("unexpected role key %q", resp.User.RoleKey)
	}
	if resp.User.UpdatedAt == nil || *resp.User.UpdatedAt != "2026-06-25T13:00:00Z" {
		t.Fatalf("unexpected updated_at %#v", resp.User.UpdatedAt)
	}
}
