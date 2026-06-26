package v1

import (
	"net/http"

	locationv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/location/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/locations/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/parse"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
)

type Handler struct {
	service      *application.Service
	defaultLimit int32
}

func NewHandler(service *application.Service, defaultLimit int32) *Handler {
	return &Handler{service: service, defaultLimit: defaultLimit}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, auth middleware.Middleware) {
	mux.Handle("GET /api/v1/locations/countries", middleware.Chain(http.HandlerFunc(h.listCountries), auth))
	mux.Handle("GET /api/v1/locations/admin-areas", middleware.Chain(http.HandlerFunc(h.listAdminAreas), auth))
	mux.Handle("GET /api/v1/locations/localities", middleware.Chain(http.HandlerFunc(h.listLocalities), auth))
	mux.Handle("GET /api/v1/locations/settlements", middleware.Chain(http.HandlerFunc(h.listSettlements), auth))
	mux.Handle("GET /api/v1/locations/postal-codes/search", middleware.Chain(http.HandlerFunc(h.searchPostalCodes), auth))
	mux.Handle("GET /api/v1/locations/address-suggestions", middleware.Chain(http.HandlerFunc(h.suggestAddress), auth))
}

func (h *Handler) listCountries(w http.ResponseWriter, r *http.Request) {
	limit, err := parse.QueryInt(r, "limit", h.defaultLimit)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.ListCountries(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &locationv1.ListCountriesRequest{
		Query: parse.OptionalString(parse.QueryString(r, "query")),
		Limit: limit,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoCountries(resp.GetCountries()), nil)
}

func (h *Handler) listAdminAreas(w http.ResponseWriter, r *http.Request) {
	limit, err := parse.QueryInt(r, "limit", h.defaultLimit)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	adminAreaType, err := parse.AdminAreaType(parse.QueryString(r, "admin_area_type"))
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.ListAdminAreas(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &locationv1.ListAdminAreasRequest{
		CountryCode:   parse.QueryString(r, "country_code"),
		ParentCode:    parse.OptionalString(parse.QueryString(r, "parent_code")),
		Query:         parse.OptionalString(parse.QueryString(r, "query")),
		Limit:         limit,
		AdminAreaType: adminAreaType,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoAdminAreas(resp.GetAdminAreas()), nil)
}

func (h *Handler) listLocalities(w http.ResponseWriter, r *http.Request) {
	limit, err := parse.QueryInt(r, "limit", h.defaultLimit)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.ListLocalities(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &locationv1.ListLocalitiesRequest{
		CountryCode:   parse.QueryString(r, "country_code"),
		AdminAreaCode: parse.OptionalString(parse.QueryString(r, "admin_area_code")),
		Query:         parse.OptionalString(parse.QueryString(r, "query")),
		Limit:         limit,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoLocalities(resp.GetLocalities()), nil)
}

func (h *Handler) listSettlements(w http.ResponseWriter, r *http.Request) {
	limit, err := parse.QueryInt(r, "limit", h.defaultLimit)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.ListSettlements(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &locationv1.ListSettlementsRequest{
		CountryCode:   parse.QueryString(r, "country_code"),
		AdminAreaCode: parse.OptionalString(parse.QueryString(r, "admin_area_code")),
		LocalityCode:  parse.OptionalString(parse.QueryString(r, "locality_code")),
		PostalCode:    parse.OptionalString(parse.QueryString(r, "postal_code")),
		Query:         parse.OptionalString(parse.QueryString(r, "query")),
		Limit:         limit,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoSettlements(resp.GetSettlements()), nil)
}

func (h *Handler) searchPostalCodes(w http.ResponseWriter, r *http.Request) {
	limit, err := parse.QueryInt(r, "limit", h.defaultLimit)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.SearchPostalCodes(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &locationv1.SearchPostalCodesRequest{
		CountryCode: parse.QueryString(r, "country_code"),
		PostalCode:  parse.QueryString(r, "postal_code"),
		Limit:       limit,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoPostalCodeMatches(resp.GetPostalCodes()), nil)
}

func (h *Handler) suggestAddress(w http.ResponseWriter, r *http.Request) {
	limit, err := parse.QueryInt(r, "limit", h.defaultLimit)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.SuggestAddress(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &locationv1.SuggestAddressRequest{
		CountryCode: parse.QueryString(r, "country_code"),
		Query:       parse.QueryString(r, "query"),
		PostalCode:  parse.OptionalString(parse.QueryString(r, "postal_code")),
		Limit:       limit,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoAddressSuggestions(resp.GetSuggestions()), nil)
}
