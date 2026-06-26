package grpc

import (
	"context"

	authv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/auth/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/auth/ports"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
)

type Client struct {
	client authv1.AuthServiceClient
}

func New(client authv1.AuthServiceClient) *Client {
	return &Client{client: client}
}

func (c *Client) SignIn(ctx context.Context, req *authv1.SignInRequest) (*authv1.SignInResponse, error) {
	resp, err := c.client.SignIn(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) SignOut(ctx context.Context, req *authv1.SignOutRequest) (*authv1.SignOutResponse, error) {
	resp, err := c.client.SignOut(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	resp, err := c.client.RefreshToken(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) RequestPasswordReset(ctx context.Context, req *authv1.RequestPasswordResetRequest) (*authv1.RequestPasswordResetResponse, error) {
	resp, err := c.client.RequestPasswordReset(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) VerifyResetCode(ctx context.Context, req *authv1.VerifyResetCodeRequest) (*authv1.VerifyResetCodeResponse, error) {
	resp, err := c.client.VerifyResetCode(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) ResetPassword(ctx context.Context, req *authv1.ResetPasswordRequest) (*authv1.ResetPasswordResponse, error) {
	resp, err := c.client.ResetPassword(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) VerifyToken(ctx context.Context, token string) (*authv1.VerifyTokenResponse, error) {
	resp, err := c.client.VerifyToken(ctx, &authv1.VerifyTokenRequest{AccessToken: token})
	return resp, grpcclient.MapError(err)
}

func (c *Client) ListPermissions(ctx context.Context, req *authv1.ListPermissionsRequest) (*authv1.ListPermissionsResponse, error) {
	resp, err := c.client.ListPermissions(ctx, req)
	return resp, grpcclient.MapError(err)
}

var _ ports.Service = (*Client)(nil)
