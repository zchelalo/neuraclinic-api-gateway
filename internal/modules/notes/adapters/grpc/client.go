package grpc

import (
	"context"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/notes/ports"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
)

type Client struct {
	client recordv1.NoteServiceClient
}

func New(client recordv1.NoteServiceClient) *Client {
	return &Client{client: client}
}

func (c *Client) Create(ctx context.Context, req *recordv1.NoteServiceCreateRequest) (*recordv1.NoteServiceCreateResponse, error) {
	resp, err := c.client.Create(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) List(ctx context.Context, req *recordv1.NoteServiceListRequest) (*recordv1.NoteServiceListResponse, error) {
	resp, err := c.client.List(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) FindByID(ctx context.Context, id string) (*recordv1.NoteServiceFindByIdResponse, error) {
	resp, err := c.client.FindById(ctx, &recordv1.NoteServiceFindByIdRequest{Id: id})
	return resp, grpcclient.MapError(err)
}

func (c *Client) Update(ctx context.Context, req *recordv1.NoteServiceUpdateRequest) (*recordv1.NoteServiceUpdateResponse, error) {
	resp, err := c.client.Update(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) Delete(ctx context.Context, id string) (*recordv1.NoteServiceDeleteResponse, error) {
	resp, err := c.client.Delete(ctx, &recordv1.NoteServiceDeleteRequest{Id: id})
	return resp, grpcclient.MapError(err)
}

var _ ports.Client = (*Client)(nil)
