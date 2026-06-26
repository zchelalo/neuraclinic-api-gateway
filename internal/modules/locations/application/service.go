package application

import (
	"context"

	locationv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/location/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/locations/ports"
)

type Service struct {
	client ports.Client
}

func NewService(client ports.Client) *Service {
	return &Service{client: client}
}

func (s *Service) ListCountries(ctx context.Context, req *locationv1.ListCountriesRequest) (*locationv1.ListCountriesResponse, error) {
	return s.client.ListCountries(ctx, req)
}

func (s *Service) ListAdminAreas(ctx context.Context, req *locationv1.ListAdminAreasRequest) (*locationv1.ListAdminAreasResponse, error) {
	return s.client.ListAdminAreas(ctx, req)
}

func (s *Service) ListLocalities(ctx context.Context, req *locationv1.ListLocalitiesRequest) (*locationv1.ListLocalitiesResponse, error) {
	return s.client.ListLocalities(ctx, req)
}

func (s *Service) ListSettlements(ctx context.Context, req *locationv1.ListSettlementsRequest) (*locationv1.ListSettlementsResponse, error) {
	return s.client.ListSettlements(ctx, req)
}

func (s *Service) SearchPostalCodes(ctx context.Context, req *locationv1.SearchPostalCodesRequest) (*locationv1.SearchPostalCodesResponse, error) {
	return s.client.SearchPostalCodes(ctx, req)
}

func (s *Service) SuggestAddress(ctx context.Context, req *locationv1.SuggestAddressRequest) (*locationv1.SuggestAddressResponse, error) {
	return s.client.SuggestAddress(ctx, req)
}
