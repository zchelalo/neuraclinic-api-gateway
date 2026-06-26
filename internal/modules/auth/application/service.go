package application

import (
	"context"

	authv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/auth/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/auth/ports"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
)

type Service struct {
	client ports.Service
}

func NewService(client ports.Service) *Service {
	return &Service{client: client}
}

func (s *Service) SignIn(ctx context.Context, req *authv1.SignInRequest) (*authv1.SignInResponse, error) {
	return s.client.SignIn(ctx, req)
}

func (s *Service) SignOut(ctx context.Context, req *authv1.SignOutRequest) (*authv1.SignOutResponse, error) {
	return s.client.SignOut(ctx, req)
}

func (s *Service) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	return s.client.RefreshToken(ctx, req)
}

func (s *Service) RequestPasswordReset(ctx context.Context, req *authv1.RequestPasswordResetRequest) (*authv1.RequestPasswordResetResponse, error) {
	return s.client.RequestPasswordReset(ctx, req)
}

func (s *Service) VerifyResetCode(ctx context.Context, req *authv1.VerifyResetCodeRequest) (*authv1.VerifyResetCodeResponse, error) {
	return s.client.VerifyResetCode(ctx, req)
}

func (s *Service) ResetPassword(ctx context.Context, req *authv1.ResetPasswordRequest) (*authv1.ResetPasswordResponse, error) {
	return s.client.ResetPassword(ctx, req)
}

func (s *Service) ListPermissions(ctx context.Context) (*authv1.ListPermissionsResponse, error) {
	return s.client.ListPermissions(ctx, &authv1.ListPermissionsRequest{})
}

func (s *Service) VerifyToken(ctx context.Context, token string) (middleware.VerifiedToken, error) {
	if token == "" {
		return middleware.VerifiedToken{}, response.Unauthorized("missing_token", "missing access token", nil)
	}
	resp, err := s.client.VerifyToken(ctx, token)
	if err != nil {
		return middleware.VerifiedToken{}, err
	}
	return middleware.VerifiedToken{
		UserID:          resp.GetUserId(),
		RoleKey:         resp.GetRoleKey(),
		PsychologistID:  resp.GetPsychologistId(),
		AdminID:         resp.GetAdminId(),
		PermissionsKeys: resp.GetPermissionsKeys(),
	}, nil
}
