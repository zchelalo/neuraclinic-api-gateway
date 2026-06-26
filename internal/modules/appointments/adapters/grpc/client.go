package grpc

import (
	"context"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/appointments/ports"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
)

type Client struct {
	client recordv1.AppointmentServiceClient
}

func New(client recordv1.AppointmentServiceClient) *Client {
	return &Client{client: client}
}

func (c *Client) Create(ctx context.Context, req *recordv1.AppointmentServiceCreateRequest) (*recordv1.AppointmentServiceCreateResponse, error) {
	resp, err := c.client.Create(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) List(ctx context.Context, req *recordv1.AppointmentServiceListRequest) (*recordv1.AppointmentServiceListResponse, error) {
	resp, err := c.client.List(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) FindByID(ctx context.Context, id string) (*recordv1.AppointmentServiceFindByIdResponse, error) {
	resp, err := c.client.FindById(ctx, &recordv1.AppointmentServiceFindByIdRequest{AppointmentId: id})
	return resp, grpcclient.MapError(err)
}

func (c *Client) Reschedule(ctx context.Context, req *recordv1.AppointmentServiceRescheduleRequest) (*recordv1.AppointmentServiceRescheduleResponse, error) {
	resp, err := c.client.Reschedule(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) UpdateStatus(ctx context.Context, req *recordv1.AppointmentServiceUpdateStatusRequest) (*recordv1.AppointmentServiceUpdateStatusResponse, error) {
	resp, err := c.client.UpdateStatus(ctx, req)
	return resp, grpcclient.MapError(err)
}

var _ ports.Client = (*Client)(nil)
