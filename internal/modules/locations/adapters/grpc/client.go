package grpc

import (
	"context"

	locationv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/location/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/locations/ports"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
)

type Client struct {
	client locationv1.LocationServiceClient
}

func New(client locationv1.LocationServiceClient) *Client {
	return &Client{client: client}
}

func (c *Client) ListCountries(ctx context.Context, req *locationv1.ListCountriesRequest) (*locationv1.ListCountriesResponse, error) {
	resp, err := c.client.ListCountries(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) ListAdminAreas(ctx context.Context, req *locationv1.ListAdminAreasRequest) (*locationv1.ListAdminAreasResponse, error) {
	resp, err := c.client.ListAdminAreas(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) ListLocalities(ctx context.Context, req *locationv1.ListLocalitiesRequest) (*locationv1.ListLocalitiesResponse, error) {
	resp, err := c.client.ListLocalities(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) ListSettlements(ctx context.Context, req *locationv1.ListSettlementsRequest) (*locationv1.ListSettlementsResponse, error) {
	resp, err := c.client.ListSettlements(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) SearchPostalCodes(ctx context.Context, req *locationv1.SearchPostalCodesRequest) (*locationv1.SearchPostalCodesResponse, error) {
	resp, err := c.client.SearchPostalCodes(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *Client) SuggestAddress(ctx context.Context, req *locationv1.SuggestAddressRequest) (*locationv1.SuggestAddressResponse, error) {
	resp, err := c.client.SuggestAddress(ctx, req)
	return resp, grpcclient.MapError(err)
}

var _ ports.Client = (*Client)(nil)
