package grpc

import (
	"context"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/patients/ports"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
)

type Client struct {
	client recordv1.PatientServiceClient
}

func New(client recordv1.PatientServiceClient) *Client {
	return &Client{client: client}
}

func (c *Client) Create(ctx context.Context, req *recordv1.PatientServiceCreateRequest) (*recordv1.PatientServiceCreateResponse, error) {
	resp, err := c.client.Create(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) List(ctx context.Context, req *recordv1.PatientServiceListRequest) (*recordv1.PatientServiceListResponse, error) {
	resp, err := c.client.List(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) FindByID(ctx context.Context, id string) (*recordv1.PatientServiceFindByIdResponse, error) {
	resp, err := c.client.FindById(ctx, &recordv1.PatientServiceFindByIdRequest{Id: id})
	return resp, grpcclient.MapError(err)
}

func (c *Client) UpdateIdentification(ctx context.Context, req *recordv1.PatientServiceUpdateIdentificationDataRequest) (*recordv1.PatientServiceUpdateIdentificationDataResponse, error) {
	resp, err := c.client.UpdateIdentificationData(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) UpdateContact(ctx context.Context, req *recordv1.PatientServiceUpdateContactDetailsRequest) (*recordv1.PatientServiceUpdateContactDetailsResponse, error) {
	resp, err := c.client.UpdateContactDetails(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) UpdateAddress(ctx context.Context, req *recordv1.PatientServiceUpdateAddressRequest) (*recordv1.PatientServiceUpdateAddressResponse, error) {
	resp, err := c.client.UpdateAddress(ctx, req)
	return resp, grpcclient.MapError(err)
}

var _ ports.Client = (*Client)(nil)
