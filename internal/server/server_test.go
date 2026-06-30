package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	authv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/auth/v1"
	locationv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/location/v1"
	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
	userv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/user/v1"
	authapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/auth/application"
	authports "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/auth/ports"
	locationsapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/locations/application"
	locationports "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/locations/ports"
	usersapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/users/application"
	userports "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/users/ports"
	"go.uber.org/zap"
)

func TestServerAuthAndProtectedRoutes(t *testing.T) {
	srv := New(Config{
		Port:                 8000,
		ServiceName:          "gateway-test",
		AccessCookieName:     "access_token",
		RefreshCookieName:    "refresh_token",
		InternalServiceToken: "internal-token",
	}, zap.NewNop(), Dependencies{
		Auth:      authapp.NewService(&fakeAuthPort{}),
		Users:     usersapp.NewService(&fakeUsersPort{}),
		Locations: locationsapp.NewService(&fakeLocationsPort{}),
	})

	t.Run("auth me via cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
		req.AddCookie(&http.Cookie{Name: "access_token", Value: "allow"})
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("auth me via bearer mobile", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
		req.Header.Set("X-Auth-Mode", "mobile")
		req.Header.Set("Authorization", "Bearer allow")
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	})

	t.Run("protected users without token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})

	t.Run("protected users without permission", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		req.AddCookie(&http.Cookie{Name: "access_token", Value: "basic"})
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403, got %d", rec.Code)
		}
	})

	t.Run("locations with auth", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/locations/countries", nil)
		req.AddCookie(&http.Cookie{Name: "access_token", Value: "basic"})
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		var body map[string]any
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["data"] == nil {
			t.Fatal("expected data in locations response")
		}
	})

	t.Run("catalogs with auth", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/catalogs/sexes", nil)
		req.AddCookie(&http.Cookie{Name: "access_token", Value: "basic"})
		req.Header.Set("Accept-Language", "es")
		rec := httptest.NewRecorder()
		srv.Handler().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		var body map[string]any
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body["data"] == nil {
			t.Fatal("expected data in catalogs response")
		}
	})
}

type fakeAuthPort struct{}

func (f *fakeAuthPort) SignIn(context.Context, *authv1.SignInRequest) (*authv1.SignInResponse, error) {
	return nil, nil
}
func (f *fakeAuthPort) SignOut(context.Context, *authv1.SignOutRequest) (*authv1.SignOutResponse, error) {
	return nil, nil
}
func (f *fakeAuthPort) RefreshToken(context.Context, *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	return nil, nil
}
func (f *fakeAuthPort) RequestPasswordReset(context.Context, *authv1.RequestPasswordResetRequest) (*authv1.RequestPasswordResetResponse, error) {
	return nil, nil
}
func (f *fakeAuthPort) VerifyResetCode(context.Context, *authv1.VerifyResetCodeRequest) (*authv1.VerifyResetCodeResponse, error) {
	return nil, nil
}
func (f *fakeAuthPort) ResetPassword(context.Context, *authv1.ResetPasswordRequest) (*authv1.ResetPasswordResponse, error) {
	return nil, nil
}
func (f *fakeAuthPort) VerifyToken(_ context.Context, token string) (*authv1.VerifyTokenResponse, error) {
	if token == "allow" {
		return &authv1.VerifyTokenResponse{
			UserId:          "user-1",
			RoleKey:         sharedv1.RoleKey_ROLE_KEY_ADMIN,
			PermissionsKeys: []string{"PERMISSION_KEY_USER_VIEW"},
		}, nil
	}
	return &authv1.VerifyTokenResponse{
		UserId:          "user-1",
		RoleKey:         sharedv1.RoleKey_ROLE_KEY_PSYCHOLOGIST,
		PermissionsKeys: nil,
	}, nil
}
func (f *fakeAuthPort) ListPermissions(context.Context, *authv1.ListPermissionsRequest) (*authv1.ListPermissionsResponse, error) {
	return nil, nil
}

type fakeUsersPort struct{}

func (f *fakeUsersPort) Create(context.Context, *userv1.UserServiceCreateRequest) (*userv1.UserServiceCreateResponse, error) {
	return nil, nil
}
func (f *fakeUsersPort) FindByID(context.Context, string) (*userv1.UserServiceFindByIdResponse, error) {
	return nil, nil
}
func (f *fakeUsersPort) FindByEmail(context.Context, string) (*userv1.UserServiceFindByEmailResponse, error) {
	return nil, nil
}
func (f *fakeUsersPort) List(context.Context, *userv1.UserServiceListRequest) (*userv1.UserServiceListResponse, error) {
	return &userv1.UserServiceListResponse{}, nil
}
func (f *fakeUsersPort) UpdatePassword(context.Context, string, string) (*userv1.UserServiceUpdatePasswordResponse, error) {
	return &userv1.UserServiceUpdatePasswordResponse{}, nil
}
func (f *fakeUsersPort) Delete(context.Context, string) (*userv1.UserServiceDeleteResponse, error) {
	return &userv1.UserServiceDeleteResponse{}, nil
}

type fakeLocationsPort struct{}

func (f *fakeLocationsPort) ListCountries(context.Context, *locationv1.ListCountriesRequest) (*locationv1.ListCountriesResponse, error) {
	return &locationv1.ListCountriesResponse{
		Countries: []*locationv1.Country{{CountryCode: "MX", Name: "Mexico"}},
	}, nil
}
func (f *fakeLocationsPort) ListAdminAreas(context.Context, *locationv1.ListAdminAreasRequest) (*locationv1.ListAdminAreasResponse, error) {
	return nil, nil
}
func (f *fakeLocationsPort) ListLocalities(context.Context, *locationv1.ListLocalitiesRequest) (*locationv1.ListLocalitiesResponse, error) {
	return nil, nil
}
func (f *fakeLocationsPort) ListSettlements(context.Context, *locationv1.ListSettlementsRequest) (*locationv1.ListSettlementsResponse, error) {
	return nil, nil
}
func (f *fakeLocationsPort) SearchPostalCodes(context.Context, *locationv1.SearchPostalCodesRequest) (*locationv1.SearchPostalCodesResponse, error) {
	return nil, nil
}
func (f *fakeLocationsPort) SuggestAddress(context.Context, *locationv1.SuggestAddressRequest) (*locationv1.SuggestAddressResponse, error) {
	return nil, nil
}

var _ authports.Service = (*fakeAuthPort)(nil)
var _ userports.Client = (*fakeUsersPort)(nil)
var _ locationports.Client = (*fakeLocationsPort)(nil)
