package v1

import (
	"net/http"
	"time"

	authv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/auth/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/auth/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/server/httpx"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/httpdto"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
)

type CookieConfig struct {
	Domain      string
	Secure      bool
	AccessName  string
	RefreshName string
}

type Handler struct {
	service *application.Service
	cookies CookieConfig
}

func NewHandler(service *application.Service, cookies CookieConfig) *Handler {
	return &Handler{service: service, cookies: cookies}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, authMiddleware middleware.Middleware) {
	mux.Handle("POST /api/v1/auth/sign-in", http.HandlerFunc(h.signIn))
	mux.Handle("POST /api/v1/auth/sign-out", http.HandlerFunc(h.signOut))
	mux.Handle("POST /api/v1/auth/refresh-token", http.HandlerFunc(h.refreshToken))
	mux.Handle("GET /api/v1/auth/me", middleware.Chain(http.HandlerFunc(h.me), authMiddleware))
	mux.Handle("POST /api/v1/auth/request-password-reset", http.HandlerFunc(h.requestPasswordReset))
	mux.Handle("POST /api/v1/auth/verify-reset-code", http.HandlerFunc(h.verifyResetCode))
	mux.Handle("POST /api/v1/auth/reset-password", http.HandlerFunc(h.resetPassword))
	mux.Handle("GET /api/v1/auth/permissions", middleware.Chain(http.HandlerFunc(h.listPermissions), authMiddleware))
}

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var body signInRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}

	resp, err := h.service.SignIn(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{}), &authv1.SignInRequest{
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	if middleware.IsMobile(r) {
		response.Write(w, r, http.StatusOK, fromProtoSignInResponse(resp), nil)
		return
	}

	h.setCookie(w, h.cookies.AccessName, resp.GetAccessToken(), resp.GetAccessTokenExpiry().AsTime())
	h.setCookie(w, h.cookies.RefreshName, resp.GetRefreshToken(), resp.GetRefreshTokenExpiry().AsTime())
	response.Write(w, r, http.StatusOK, webSignInResponse{
		SignedIn:           true,
		AccessTokenExpiry:  httpdto.Timestamp(resp.GetAccessTokenExpiry()),
		RefreshTokenExpiry: httpdto.Timestamp(resp.GetRefreshTokenExpiry()),
	}, nil)
}

type signOutRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) signOut(w http.ResponseWriter, r *http.Request) {
	var body signOutRequest
	if err := httpx.DecodeBodyAllowEmpty(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}

	accessToken := body.AccessToken
	refreshToken := body.RefreshToken
	if middleware.IsMobile(r) && accessToken == "" {
		accessToken = middleware.BearerToken(r)
	}
	if !middleware.IsMobile(r) {
		if cookie, err := r.Cookie(h.cookies.AccessName); err == nil && accessToken == "" {
			accessToken = cookie.Value
		}
		if cookie, err := r.Cookie(h.cookies.RefreshName); err == nil && refreshToken == "" {
			refreshToken = cookie.Value
		}
	}

	resp, err := h.service.SignOut(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{}), &authv1.SignOutRequest{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	if !middleware.IsMobile(r) {
		h.clearCookie(w, h.cookies.AccessName)
		h.clearCookie(w, h.cookies.RefreshName)
	}
	response.Write(w, r, http.StatusOK, httpdto.Operation(resp.GetOperation()), nil)
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) refreshToken(w http.ResponseWriter, r *http.Request) {
	var body refreshTokenRequest
	if err := httpx.DecodeBodyAllowEmpty(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}

	refreshToken := body.RefreshToken
	if !middleware.IsMobile(r) && refreshToken == "" {
		if cookie, err := r.Cookie(h.cookies.RefreshName); err == nil {
			refreshToken = cookie.Value
		}
	}
	if refreshToken == "" {
		response.WriteError(w, r, response.Unauthorized("missing_refresh_token", "missing refresh token", nil))
		return
	}

	resp, err := h.service.RefreshToken(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{}), &authv1.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}

	if middleware.IsMobile(r) {
		response.Write(w, r, http.StatusOK, fromProtoRefreshTokenResponse(resp), nil)
		return
	}

	h.setCookie(w, h.cookies.AccessName, resp.GetAccessToken(), resp.GetAccessTokenExpiry().AsTime())
	if token := resp.GetRefreshToken(); token != "" && resp.GetRefreshTokenExpiry() != nil {
		h.setCookie(w, h.cookies.RefreshName, token, resp.GetRefreshTokenExpiry().AsTime())
	}
	response.Write(w, r, http.StatusOK, webRefreshResponse{
		Refreshed:           true,
		AccessTokenExpiry:   httpdto.Timestamp(resp.GetAccessTokenExpiry()),
		RefreshTokenExpiry:  httpdto.Timestamp(resp.GetRefreshTokenExpiry()),
		RefreshTokenUpdated: resp.GetRefreshToken() != "",
	}, nil)
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.Auth(r.Context())
	if !ok {
		response.WriteError(w, r, response.Unauthorized("missing_auth_context", "missing auth context", nil))
		return
	}
	response.Write(w, r, http.StatusOK, fromAuthContext(auth), nil)
}

type emailRequest struct {
	Email string `json:"email"`
}

func (h *Handler) requestPasswordReset(w http.ResponseWriter, r *http.Request) {
	var body emailRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.RequestPasswordReset(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{}), &authv1.RequestPasswordResetRequest{
		Email: body.Email,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, httpdto.Operation(resp.GetOperation()), nil)
}

type verifyResetCodeRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func (h *Handler) verifyResetCode(w http.ResponseWriter, r *http.Request) {
	var body verifyResetCodeRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.VerifyResetCode(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{}), &authv1.VerifyResetCodeRequest{
		Email: body.Email,
		Otp:   body.OTP,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoVerifyResetCodeResponse(resp), nil)
}

type resetPasswordRequest struct {
	ResetToken  string `json:"reset_token"`
	NewPassword string `json:"new_password"`
}

func (h *Handler) resetPassword(w http.ResponseWriter, r *http.Request) {
	var body resetPasswordRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.ResetPassword(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{}), &authv1.ResetPasswordRequest{
		ResetToken:  body.ResetToken,
		NewPassword: body.NewPassword,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, httpdto.Operation(resp.GetOperation()), nil)
}

func (h *Handler) listPermissions(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.ListPermissions(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}))
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoPermissions(resp.GetPermissions()), nil)
}

func (h *Handler) setCookie(w http.ResponseWriter, name, value string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   h.cookies.Domain,
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.cookies.Secure,
	})
}

func (h *Handler) clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Domain:   h.cookies.Domain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.cookies.Secure,
	})
}
