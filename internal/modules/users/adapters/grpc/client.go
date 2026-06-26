package grpc

import (
	"context"

	userv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/user/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/users/ports"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
)

type Client struct {
	client userv1.UserServiceClient
}

func New(client userv1.UserServiceClient) *Client {
	return &Client{client: client}
}

func (c *Client) Create(ctx context.Context, req *userv1.UserServiceCreateRequest) (*userv1.UserServiceCreateResponse, error) {
	resp, err := c.client.Create(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) FindByID(ctx context.Context, id string) (*userv1.UserServiceFindByIdResponse, error) {
	resp, err := c.client.FindById(ctx, &userv1.UserServiceFindByIdRequest{Id: id})
	return resp, grpcclient.MapError(err)
}

func (c *Client) FindByEmail(ctx context.Context, email string) (*userv1.UserServiceFindByEmailResponse, error) {
	resp, err := c.client.FindByEmail(ctx, &userv1.UserServiceFindByEmailRequest{Email: email})
	return resp, grpcclient.MapError(err)
}

func (c *Client) List(ctx context.Context, req *userv1.UserServiceListRequest) (*userv1.UserServiceListResponse, error) {
	resp, err := c.client.List(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) UpdatePassword(ctx context.Context, id, newPassword string) (*userv1.UserServiceUpdatePasswordResponse, error) {
	resp, err := c.client.UpdatePassword(ctx, &userv1.UserServiceUpdatePasswordRequest{
		Id:          id,
		NewPassword: newPassword,
	})
	return resp, grpcclient.MapError(err)
}

func (c *Client) Delete(ctx context.Context, id string) (*userv1.UserServiceDeleteResponse, error) {
	resp, err := c.client.Delete(ctx, &userv1.UserServiceDeleteRequest{Id: id})
	return resp, grpcclient.MapError(err)
}

var _ ports.Client = (*Client)(nil)
