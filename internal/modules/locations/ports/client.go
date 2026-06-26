package ports

import (
	"context"

	locationv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/location/v1"
)

type Client interface {
	ListCountries(ctx context.Context, req *locationv1.ListCountriesRequest) (*locationv1.ListCountriesResponse, error)
	ListAdminAreas(ctx context.Context, req *locationv1.ListAdminAreasRequest) (*locationv1.ListAdminAreasResponse, error)
	ListLocalities(ctx context.Context, req *locationv1.ListLocalitiesRequest) (*locationv1.ListLocalitiesResponse, error)
	ListSettlements(ctx context.Context, req *locationv1.ListSettlementsRequest) (*locationv1.ListSettlementsResponse, error)
	SearchPostalCodes(ctx context.Context, req *locationv1.SearchPostalCodesRequest) (*locationv1.SearchPostalCodesResponse, error)
	SuggestAddress(ctx context.Context, req *locationv1.SuggestAddressRequest) (*locationv1.SuggestAddressResponse, error)
}
