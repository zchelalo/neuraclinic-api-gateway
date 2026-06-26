package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	locationv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/location/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/locations/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/locations/ports"
)

func TestListCountriesContract(t *testing.T) {
	handler := NewHandler(application.NewService(&fakeLocationsClient{
		listCountriesResp: &locationv1.ListCountriesResponse{
			Countries: []*locationv1.Country{
				{CountryCode: "MX", Name: "Mexico", Label: "Mexico", Source: "sepomex", SourceVersion: "2026"},
			},
		},
	}), 10)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, passthroughLocationsMiddleware())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/locations/countries", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	assertLocationsJSONEqual(t, rec.Body.Bytes(), `{
		"data": [{
			"country_code": "MX",
			"name": "Mexico",
			"label": "Mexico",
			"source": "sepomex",
			"source_version": "2026"
		}],
		"request_id": ""
	}`)
}

type fakeLocationsClient struct {
	listCountriesResp *locationv1.ListCountriesResponse
}

func (f *fakeLocationsClient) ListCountries(context.Context, *locationv1.ListCountriesRequest) (*locationv1.ListCountriesResponse, error) {
	return f.listCountriesResp, nil
}
func (f *fakeLocationsClient) ListAdminAreas(context.Context, *locationv1.ListAdminAreasRequest) (*locationv1.ListAdminAreasResponse, error) {
	return nil, nil
}
func (f *fakeLocationsClient) ListLocalities(context.Context, *locationv1.ListLocalitiesRequest) (*locationv1.ListLocalitiesResponse, error) {
	return nil, nil
}
func (f *fakeLocationsClient) ListSettlements(context.Context, *locationv1.ListSettlementsRequest) (*locationv1.ListSettlementsResponse, error) {
	return nil, nil
}
func (f *fakeLocationsClient) SearchPostalCodes(context.Context, *locationv1.SearchPostalCodesRequest) (*locationv1.SearchPostalCodesResponse, error) {
	return nil, nil
}
func (f *fakeLocationsClient) SuggestAddress(context.Context, *locationv1.SuggestAddressRequest) (*locationv1.SuggestAddressResponse, error) {
	return nil, nil
}

var _ ports.Client = (*fakeLocationsClient)(nil)

func passthroughLocationsMiddleware() middleware.Middleware {
	return func(next http.Handler) http.Handler { return next }
}

func assertLocationsJSONEqual(t *testing.T, got []byte, expected string) {
	t.Helper()
	var gotValue any
	var expectedValue any
	if err := json.Unmarshal(got, &gotValue); err != nil {
		t.Fatalf("unmarshal got: %v", err)
	}
	if err := json.Unmarshal([]byte(expected), &expectedValue); err != nil {
		t.Fatalf("unmarshal expected: %v", err)
	}
	if !reflect.DeepEqual(gotValue, expectedValue) {
		t.Fatalf("unexpected json\n got: %s\nwant: %s", got, expected)
	}
}
