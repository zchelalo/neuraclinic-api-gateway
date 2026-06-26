package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
	userv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/user/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/users/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/users/ports"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestUsersContractResponses(t *testing.T) {
	now := timestamppb.New(time.Date(2026, 6, 25, 13, 0, 0, 0, time.UTC))
	client := &fakeUsersClient{
		findByIDResp: &userv1.UserServiceFindByIdResponse{
			User: &userv1.User{
				Id:        "user-1",
				Email:     "psy@example.com",
				RoleKey:   sharedv1.RoleKey_ROLE_KEY_PSYCHOLOGIST,
				CreatedAt: now,
			},
			Profile: &userv1.UserServiceFindByIdResponse_Psychologist{
				Psychologist: &userv1.PsychologistProfile{
					Id:            "psy-1",
					FirstName:     "Ana",
					MiddleName:    strPtr("Maria"),
					FirstLastName: "Lopez",
					UpdatedAt:     now,
				},
			},
		},
		listResp: &userv1.UserServiceListResponse{
			Users: []*userv1.User{
				{Id: "user-1", Email: "admin@example.com", RoleKey: sharedv1.RoleKey_ROLE_KEY_ADMIN, CreatedAt: now},
			},
			Meta: &sharedv1.CursorMeta{
				NextCursor: strPtr("next-1"),
				Limit:      int32Ptr(10),
			},
		},
		updatePasswordResp: &userv1.UserServiceUpdatePasswordResponse{Operation: operation("updated")},
	}
	handler := NewHandler(application.NewService(client), "internal-token", 10)
	auth := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := middleware.WithAuth(r.Context(), middleware.AuthContext{UserID: "self-id", PermissionsKeys: []string{permUserView}})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, auth, func(...string) middleware.Middleware { return passthroughMiddleware() })

	t.Run("find by id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/user-1", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assertUsersJSONEqual(t, rec.Body.Bytes(), `{
			"data": {
				"user": {
					"id": "user-1",
					"email": "psy@example.com",
					"role_key": "ROLE_KEY_PSYCHOLOGIST",
					"created_at": "2026-06-25T13:00:00Z"
				},
				"psychologist": {
					"id": "psy-1",
					"first_name": "Ana",
					"middle_name": "Maria",
					"first_last_name": "Lopez",
					"updated_at": "2026-06-25T13:00:00Z"
				}
			},
			"request_id": ""
		}`)
	})

	t.Run("list with meta", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assertUsersJSONEqual(t, rec.Body.Bytes(), `{
			"data": [{
				"id": "user-1",
				"email": "admin@example.com",
				"role_key": "ROLE_KEY_ADMIN",
				"created_at": "2026-06-25T13:00:00Z"
			}],
			"meta": {
				"next_cursor": "next-1",
				"limit": 10
			},
			"request_id": ""
		}`)
	})
}

func TestUpdatePasswordGuard(t *testing.T) {
	client := &fakeUsersClient{updatePasswordResp: &userv1.UserServiceUpdatePasswordResponse{}}
	handler := NewHandler(application.NewService(client), "internal-token", 10)
	auth := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := middleware.WithAuth(r.Context(), middleware.AuthContext{UserID: "self-id"})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, auth, middleware.PermissionsMiddleware)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/self-id/password", strings.NewReader(`{"new_password":"new-pass"}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected self update to pass, got status %d", rec.Code)
	}
	if client.lastUpdateID != "self-id" {
		t.Fatalf("expected update password for self-id, got %q", client.lastUpdateID)
	}

	req = httptest.NewRequest(http.MethodPatch, "/api/v1/users/other-id/password", strings.NewReader(`{"new_password":"new-pass"}`))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden for other user without permission, got %d", rec.Code)
	}
}

type fakeUsersClient struct {
	lastUpdateID       string
	findByIDResp       *userv1.UserServiceFindByIdResponse
	listResp           *userv1.UserServiceListResponse
	updatePasswordResp *userv1.UserServiceUpdatePasswordResponse
}

func (f *fakeUsersClient) Create(context.Context, *userv1.UserServiceCreateRequest) (*userv1.UserServiceCreateResponse, error) {
	return nil, nil
}

func (f *fakeUsersClient) FindByID(context.Context, string) (*userv1.UserServiceFindByIdResponse, error) {
	return f.findByIDResp, nil
}

func (f *fakeUsersClient) FindByEmail(context.Context, string) (*userv1.UserServiceFindByEmailResponse, error) {
	return nil, nil
}

func (f *fakeUsersClient) List(context.Context, *userv1.UserServiceListRequest) (*userv1.UserServiceListResponse, error) {
	return f.listResp, nil
}

func (f *fakeUsersClient) UpdatePassword(_ context.Context, id, _ string) (*userv1.UserServiceUpdatePasswordResponse, error) {
	f.lastUpdateID = id
	return f.updatePasswordResp, nil
}

func (f *fakeUsersClient) Delete(context.Context, string) (*userv1.UserServiceDeleteResponse, error) {
	return nil, nil
}

var _ ports.Client = (*fakeUsersClient)(nil)

func assertUsersJSONEqual(t *testing.T, got []byte, expected string) {
	t.Helper()

	var gotValue any
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Fatalf("unmarshal got: %v", err)
	}

	var expectedValue any
	if err := json.Unmarshal([]byte(expected), &expectedValue); err != nil {
		t.Fatalf("unmarshal expected: %v", err)
	}

	if !reflect.DeepEqual(gotValue, expectedValue) {
		t.Fatalf("unexpected json\n got: %s\nwant: %s", got, expected)
	}
}

func int32Ptr(value int32) *int32 {
	return &value
}

func strPtr(value string) *string {
	return &value
}

func operation(message string) *sharedv1.OperationResponse {
	return &sharedv1.OperationResponse{Message: message}
}

func passthroughMiddleware() middleware.Middleware {
	return func(next http.Handler) http.Handler { return next }
}
