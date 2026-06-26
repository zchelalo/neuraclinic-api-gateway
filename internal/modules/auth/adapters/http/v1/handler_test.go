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

	authv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/auth/v1"
	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/auth/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/auth/ports"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAuthContractResponses(t *testing.T) {
	accessExpiry := timestamppb.New(time.Date(2026, 6, 25, 12, 0, 0, 0, time.UTC))
	refreshExpiry := timestamppb.New(time.Date(2026, 6, 26, 12, 0, 0, 0, time.UTC))
	service := application.NewService(&fakeAuthService{
		signInResp: &authv1.SignInResponse{
			AccessToken:        "access-token",
			RefreshToken:       "refresh-token",
			AccessTokenExpiry:  accessExpiry,
			RefreshTokenExpiry: refreshExpiry,
		},
		refreshResp: &authv1.RefreshTokenResponse{
			AccessToken:        "new-access",
			RefreshToken:       strPtr("new-refresh"),
			AccessTokenExpiry:  accessExpiry,
			RefreshTokenExpiry: refreshExpiry,
		},
		verifyResetCodeResp: &authv1.VerifyResetCodeResponse{ResetToken: "reset-token"},
		listPermissionsResp: &authv1.ListPermissionsResponse{
			Permissions: []*authv1.Permission{
				{
					Id:          "perm-1",
					Key:         sharedv1.PermissionKey_PERMISSION_KEY_USER_VIEW,
					Description: "Can view users",
					CreatedAt:   accessExpiry,
				},
			},
		},
	})

	handler := NewHandler(service, CookieConfig{
		AccessName:  "access_token",
		RefreshName: "refresh_token",
	})
	mux := http.NewServeMux()
	auth := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := middleware.WithAuth(r.Context(), middleware.AuthContext{
				UserID:          "user-1",
				RoleKey:         "ROLE_KEY_ADMIN",
				PsychologistID:  "",
				AdminID:         "admin-1",
				PermissionsKeys: []string{"PERMISSION_KEY_USER_VIEW"},
				Mode:            middleware.AuthModeMobile,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	handler.RegisterRoutes(mux, auth)

	t.Run("sign in web", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", strings.NewReader(`{"email":"a@example.com","password":"secret"}`))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assertJSONEqual(t, rec.Body.Bytes(), `{
			"data": {
				"signed_in": true,
				"access_token_expiry": "2026-06-25T12:00:00Z",
				"refresh_token_expiry": "2026-06-26T12:00:00Z"
			},
			"request_id": ""
		}`)
	})

	t.Run("sign in mobile", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", strings.NewReader(`{"email":"a@example.com","password":"secret"}`))
		req.Header.Set("X-Auth-Mode", "mobile")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assertJSONEqual(t, rec.Body.Bytes(), `{
			"data": {
				"access_token": "access-token",
				"refresh_token": "refresh-token",
				"access_token_expiry": "2026-06-25T12:00:00Z",
				"refresh_token_expiry": "2026-06-26T12:00:00Z"
			},
			"request_id": ""
		}`)
	})

	t.Run("me", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assertJSONEqual(t, rec.Body.Bytes(), `{
			"data": {
				"user_id": "user-1",
				"role_key": "ROLE_KEY_ADMIN",
				"psychologist_id": "",
				"admin_id": "admin-1",
				"permissions_keys": ["PERMISSION_KEY_USER_VIEW"],
				"mode": "mobile"
			},
			"request_id": ""
		}`)
	})

	t.Run("verify reset code", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/verify-reset-code", strings.NewReader(`{"email":"a@example.com","otp":"123456"}`))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assertJSONEqual(t, rec.Body.Bytes(), `{
			"data": {
				"reset_token": "reset-token"
			},
			"request_id": ""
		}`)
	})

	t.Run("permissions", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/permissions", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assertJSONEqual(t, rec.Body.Bytes(), `{
			"data": [{
				"id": "perm-1",
				"key": "PERMISSION_KEY_USER_VIEW",
				"description": "Can view users",
				"created_at": "2026-06-25T12:00:00Z"
			}],
			"request_id": ""
		}`)
	})
}

func TestSignInAndRefreshSetCookiesAndSignOutClearsThem(t *testing.T) {
	service := application.NewService(&fakeAuthService{
		signInResp: &authv1.SignInResponse{
			AccessToken:        "access-token",
			RefreshToken:       "refresh-token",
			AccessTokenExpiry:  timestamppb.New(time.Now().Add(time.Hour)),
			RefreshTokenExpiry: timestamppb.New(time.Now().Add(24 * time.Hour)),
		},
		refreshResp: &authv1.RefreshTokenResponse{
			AccessToken:        "new-access",
			RefreshToken:       strPtr("new-refresh"),
			AccessTokenExpiry:  timestamppb.New(time.Now().Add(time.Hour)),
			RefreshTokenExpiry: timestamppb.New(time.Now().Add(24 * time.Hour)),
		},
		signOutResp: &authv1.SignOutResponse{Operation: operation("signed out")},
	})

	handler := NewHandler(service, CookieConfig{
		AccessName:  "access_token",
		RefreshName: "refresh_token",
	})
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, passthroughMiddleware())

	t.Run("sign in", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", strings.NewReader(`{"email":"a@example.com","password":"secret"}`))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		cookies := rec.Result().Cookies()
		if len(cookies) != 2 {
			t.Fatalf("expected 2 cookies, got %d", len(cookies))
		}
	})

	t.Run("refresh", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh-token", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "refresh-token"})
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		if len(rec.Result().Cookies()) == 0 {
			t.Fatal("expected cookies to be set on refresh")
		}
	})

	t.Run("sign out", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-out", nil)
		req.AddCookie(&http.Cookie{Name: "access_token", Value: "access-token"})
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "refresh-token"})
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		cookies := rec.Result().Cookies()
		if len(cookies) < 2 {
			t.Fatalf("expected cleared cookies, got %d", len(cookies))
		}
		for _, cookie := range cookies {
			if cookie.MaxAge != -1 {
				t.Fatalf("expected cookie %s to be cleared", cookie.Name)
			}
		}
	})
}

type fakeAuthService struct {
	signInResp          *authv1.SignInResponse
	refreshResp         *authv1.RefreshTokenResponse
	signOutResp         *authv1.SignOutResponse
	verifyResetCodeResp *authv1.VerifyResetCodeResponse
	listPermissionsResp *authv1.ListPermissionsResponse
}

func (f *fakeAuthService) SignIn(context.Context, *authv1.SignInRequest) (*authv1.SignInResponse, error) {
	return f.signInResp, nil
}

func (f *fakeAuthService) SignOut(context.Context, *authv1.SignOutRequest) (*authv1.SignOutResponse, error) {
	return f.signOutResp, nil
}

func (f *fakeAuthService) RefreshToken(context.Context, *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	return f.refreshResp, nil
}

func (f *fakeAuthService) RequestPasswordReset(context.Context, *authv1.RequestPasswordResetRequest) (*authv1.RequestPasswordResetResponse, error) {
	return nil, nil
}

func (f *fakeAuthService) VerifyResetCode(context.Context, *authv1.VerifyResetCodeRequest) (*authv1.VerifyResetCodeResponse, error) {
	return f.verifyResetCodeResp, nil
}

func (f *fakeAuthService) ResetPassword(context.Context, *authv1.ResetPasswordRequest) (*authv1.ResetPasswordResponse, error) {
	return nil, nil
}

func (f *fakeAuthService) VerifyToken(context.Context, string) (*authv1.VerifyTokenResponse, error) {
	return nil, nil
}

func (f *fakeAuthService) ListPermissions(context.Context, *authv1.ListPermissionsRequest) (*authv1.ListPermissionsResponse, error) {
	return f.listPermissionsResp, nil
}

var _ ports.Service = (*fakeAuthService)(nil)

func operation(message string) *sharedv1.OperationResponse {
	return &sharedv1.OperationResponse{Message: message}
}

func strPtr(value string) *string {
	return &value
}

func passthroughMiddleware() middleware.Middleware {
	return func(next http.Handler) http.Handler { return next }
}

func assertJSONEqual(t *testing.T, got []byte, expected string) {
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
