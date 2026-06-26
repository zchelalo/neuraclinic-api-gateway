package ports

import (
	"context"

	authv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/auth/v1"
)

type Service interface {
	SignIn(ctx context.Context, req *authv1.SignInRequest) (*authv1.SignInResponse, error)
	SignOut(ctx context.Context, req *authv1.SignOutRequest) (*authv1.SignOutResponse, error)
	RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error)
	RequestPasswordReset(ctx context.Context, req *authv1.RequestPasswordResetRequest) (*authv1.RequestPasswordResetResponse, error)
	VerifyResetCode(ctx context.Context, req *authv1.VerifyResetCodeRequest) (*authv1.VerifyResetCodeResponse, error)
	ResetPassword(ctx context.Context, req *authv1.ResetPasswordRequest) (*authv1.ResetPasswordResponse, error)
	VerifyToken(ctx context.Context, token string) (*authv1.VerifyTokenResponse, error)
	ListPermissions(ctx context.Context, req *authv1.ListPermissionsRequest) (*authv1.ListPermissionsResponse, error)
}
