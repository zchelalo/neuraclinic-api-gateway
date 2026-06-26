package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
)

const headerAuthMode = "X-Auth-Mode"

type TokenVerifier interface {
	VerifyToken(ctx context.Context, token string) (VerifiedToken, error)
}

type VerifiedToken struct {
	UserID          string
	RoleKey         sharedv1.RoleKey
	PsychologistID  string
	AdminID         string
	PermissionsKeys []string
}

func AuthMiddleware(verifier TokenVerifier, accessCookieName string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, mode, err := ExtractToken(r, accessCookieName)
			if err != nil {
				response.WriteError(w, r, err)
				return
			}

			verified, err := verifier.VerifyToken(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{}), token)
			if err != nil {
				response.WriteError(w, r, err)
				return
			}

			ctx := WithAuth(r.Context(), AuthContext{
				Token:           token,
				Mode:            mode,
				UserID:          verified.UserID,
				RoleKey:         verified.RoleKey.String(),
				PsychologistID:  verified.PsychologistID,
				AdminID:         verified.AdminID,
				PermissionsKeys: verified.PermissionsKeys,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ExtractToken(r *http.Request, accessCookieName string) (string, AuthMode, error) {
	mode := AuthModeWeb
	if strings.EqualFold(strings.TrimSpace(r.Header.Get(headerAuthMode)), string(AuthModeMobile)) {
		mode = AuthModeMobile
	}

	if mode == AuthModeMobile {
		authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
		if authHeader == "" {
			return "", mode, response.Unauthorized("missing_token", "missing bearer token", nil)
		}
		if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			return "", mode, response.Unauthorized("invalid_token", "authorization header must use bearer token", nil)
		}
		token := strings.TrimSpace(authHeader[len("Bearer "):])
		if token == "" {
			return "", mode, response.Unauthorized("missing_token", "missing bearer token", nil)
		}
		return token, mode, nil
	}

	cookie, err := r.Cookie(accessCookieName)
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		return "", mode, response.Unauthorized("missing_token", "missing access token cookie", err)
	}
	return strings.TrimSpace(cookie.Value), mode, nil
}

func IsMobile(r *http.Request) bool {
	return strings.EqualFold(strings.TrimSpace(r.Header.Get(headerAuthMode)), string(AuthModeMobile))
}

func BearerToken(r *http.Request) string {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return ""
	}
	return strings.TrimSpace(authHeader[len("Bearer "):])
}

func RequireAuth(r *http.Request) (AuthContext, error) {
	auth, ok := Auth(r.Context())
	if !ok {
		return AuthContext{}, response.Unauthorized("missing_auth_context", "missing auth context", nil)
	}
	if auth.UserID == "" {
		return AuthContext{}, response.Unauthorized("invalid_auth_context", "invalid auth context", nil)
	}
	return auth, nil
}

func RequireCurrentUserOrPermission(r *http.Request, userID string, permission string) error {
	auth, err := RequireAuth(r)
	if err != nil {
		return err
	}
	if auth.UserID == userID {
		return nil
	}
	for _, current := range auth.PermissionsKeys {
		if current == permission {
			return nil
		}
	}
	return response.Forbidden("permission_denied", fmt.Sprintf("missing required permission %s", permission), nil)
}
