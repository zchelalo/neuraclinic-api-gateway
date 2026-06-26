package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestExtractTokenWebCookie(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(cookie("access_token", "cookie-token"))

	token, mode, err := ExtractToken(req, "access_token")
	if err != nil {
		t.Fatalf("ExtractToken returned error: %v", err)
	}
	if token != "cookie-token" {
		t.Fatalf("expected cookie token, got %q", token)
	}
	if mode != AuthModeWeb {
		t.Fatalf("expected web mode, got %q", mode)
	}
}

func TestExtractTokenMobileBearer(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Auth-Mode", "mobile")
	req.Header.Set("Authorization", "Bearer mobile-token")

	token, mode, err := ExtractToken(req, "access_token")
	if err != nil {
		t.Fatalf("ExtractToken returned error: %v", err)
	}
	if token != "mobile-token" {
		t.Fatalf("expected bearer token, got %q", token)
	}
	if mode != AuthModeMobile {
		t.Fatalf("expected mobile mode, got %q", mode)
	}
}

func TestRequireCurrentUserOrPermission(t *testing.T) {
	req := httptest.NewRequest("PATCH", "/", nil)
	req = req.WithContext(WithAuth(req.Context(), AuthContext{
		UserID:          "user-1",
		PermissionsKeys: []string{"PERMISSION_KEY_USER_EDIT"},
	}))

	if err := RequireCurrentUserOrPermission(req, "user-1", "PERMISSION_KEY_USER_EDIT"); err != nil {
		t.Fatalf("expected self update to pass, got %v", err)
	}
	if err := RequireCurrentUserOrPermission(req, "user-2", "PERMISSION_KEY_USER_EDIT"); err != nil {
		t.Fatalf("expected permission-based update to pass, got %v", err)
	}

	req = req.WithContext(WithAuth(req.Context(), AuthContext{
		UserID:          "user-1",
		PermissionsKeys: nil,
	}))
	if err := RequireCurrentUserOrPermission(req, "user-2", "PERMISSION_KEY_USER_EDIT"); err == nil {
		t.Fatal("expected forbidden error when acting on another user without permission")
	}
}

func TestAuthMiddlewareForwardsResolvedLanguageToVerifier(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(cookie("access_token", "cookie-token"))
	req.Header.Set("Accept-Language", "es-MX,en;q=0.8")

	verifier := &fakeTokenVerifier{}
	handler := AuthMiddleware(verifier, "access_token")(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))

	handler.ServeHTTP(httptest.NewRecorder(), req)

	if verifier.token != "cookie-token" {
		t.Fatalf("VerifyToken token = %q, want %q", verifier.token, "cookie-token")
	}
	if verifier.acceptLanguage != "es" {
		t.Fatalf("VerifyToken accept-language = %q, want %q", verifier.acceptLanguage, "es")
	}
}

func cookie(name, value string) *http.Cookie {
	return &http.Cookie{Name: name, Value: value}
}

type fakeTokenVerifier struct {
	token          string
	acceptLanguage string
}

func (f *fakeTokenVerifier) VerifyToken(ctx context.Context, token string) (VerifiedToken, error) {
	f.token = token
	md, _ := metadata.FromOutgoingContext(ctx)
	f.acceptLanguage = first(md.Get("accept-language"))
	return VerifiedToken{UserID: "user-1"}, nil
}

func first(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
