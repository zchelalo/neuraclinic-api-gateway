package grpc

import (
	"context"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/familiograms/ports"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
)

type Client struct {
	client recordv1.FamiliogramServiceClient
}

func New(client recordv1.FamiliogramServiceClient) *Client {
	return &Client{client: client}
}

func (c *Client) FindByPatientID(ctx context.Context, patientID string) (*recordv1.FamiliogramServiceFindByPatientIdResponse, error) {
	resp, err := c.client.FindByPatientId(ctx, &recordv1.FamiliogramServiceFindByPatientIdRequest{PatientId: patientID})
	return resp, grpcclient.MapError(err)
}

func (c *Client) Update(ctx context.Context, req *recordv1.FamiliogramServiceUpdateRequest) (*recordv1.FamiliogramServiceUpdateResponse, error) {
	resp, err := c.client.Update(ctx, req)
	return resp, grpcclient.MapError(err)
}

var _ ports.Client = (*Client)(nil)
