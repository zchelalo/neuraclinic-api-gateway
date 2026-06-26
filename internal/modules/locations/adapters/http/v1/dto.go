package v1

import (
	locationv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/location/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/httpdto"
)

type locationComponentsResponse struct {
	CountryCode    string  `json:"country_code,omitempty"`
	CountryName    string  `json:"country_name,omitempty"`
	AdminAreaCode  *string `json:"admin_area_code,omitempty"`
	AdminAreaName  *string `json:"admin_area_name,omitempty"`
	LocalityCode   *string `json:"locality_code,omitempty"`
	LocalityName   *string `json:"locality_name,omitempty"`
	PostalCode     *string `json:"postal_code,omitempty"`
	SettlementName *string `json:"settlement_name,omitempty"`
	SettlementType *string `json:"settlement_type,omitempty"`
	StreetName     *string `json:"street_name,omitempty"`
}

type countryResponse struct {
	CountryCode   string `json:"country_code,omitempty"`
	Name          string `json:"name,omitempty"`
	Label         string `json:"label,omitempty"`
	Source        string `json:"source,omitempty"`
	SourceVersion string `json:"source_version,omitempty"`
}

type adminAreaResponse struct {
	Id            string  `json:"id,omitempty"`
	CountryCode   string  `json:"country_code,omitempty"`
	Code          string  `json:"code,omitempty"`
	Name          string  `json:"name,omitempty"`
	ParentCode    *string `json:"parent_code,omitempty"`
	Label         string  `json:"label,omitempty"`
	Source        string  `json:"source,omitempty"`
	SourceVersion string  `json:"source_version,omitempty"`
	AdminAreaType string  `json:"admin_area_type,omitempty"`
}

type localityResponse struct {
	Id            string `json:"id,omitempty"`
	CountryCode   string `json:"country_code,omitempty"`
	AdminAreaCode string `json:"admin_area_code,omitempty"`
	Code          string `json:"code,omitempty"`
	Name          string `json:"name,omitempty"`
	Type          string `json:"type,omitempty"`
	Label         string `json:"label,omitempty"`
	Source        string `json:"source,omitempty"`
	SourceVersion string `json:"source_version,omitempty"`
}

type settlementResponse struct {
	Id            string  `json:"id,omitempty"`
	CountryCode   string  `json:"country_code,omitempty"`
	AdminAreaCode string  `json:"admin_area_code,omitempty"`
	LocalityCode  *string `json:"locality_code,omitempty"`
	PostalCode    *string `json:"postal_code,omitempty"`
	Name          string  `json:"name,omitempty"`
	Type          string  `json:"type,omitempty"`
	Label         string  `json:"label,omitempty"`
	Source        string  `json:"source,omitempty"`
	SourceVersion string  `json:"source_version,omitempty"`
}

type postalCodeMatchResponse struct {
	PostalCode    string                      `json:"postal_code,omitempty"`
	Label         string                      `json:"label,omitempty"`
	Components    *locationComponentsResponse `json:"components,omitempty"`
	Source        string                      `json:"source,omitempty"`
	SourceVersion string                      `json:"source_version,omitempty"`
	Score         float64                     `json:"score,omitempty"`
}

type addressSuggestionResponse struct {
	Label         string                      `json:"label,omitempty"`
	Components    *locationComponentsResponse `json:"components,omitempty"`
	Source        string                      `json:"source,omitempty"`
	SourceVersion string                      `json:"source_version,omitempty"`
	Score         float64                     `json:"score,omitempty"`
}

func fromProtoLocationComponents(value *locationv1.LocationComponents) *locationComponentsResponse {
	if value == nil {
		return nil
	}
	return &locationComponentsResponse{
		CountryCode:    value.GetCountryCode(),
		CountryName:    value.GetCountryName(),
		AdminAreaCode:  value.AdminAreaCode,
		AdminAreaName:  value.AdminAreaName,
		LocalityCode:   value.LocalityCode,
		LocalityName:   value.LocalityName,
		PostalCode:     value.PostalCode,
		SettlementName: value.SettlementName,
		SettlementType: value.SettlementType,
		StreetName:     value.StreetName,
	}
}

func fromProtoCountries(values []*locationv1.Country) []countryResponse {
	if values == nil {
		return nil
	}
	result := make([]countryResponse, 0, len(values))
	for _, value := range values {
		result = append(result, countryResponse{
			CountryCode:   value.GetCountryCode(),
			Name:          value.GetName(),
			Label:         value.GetLabel(),
			Source:        value.GetSource(),
			SourceVersion: value.GetSourceVersion(),
		})
	}
	return result
}

func fromProtoAdminAreas(values []*locationv1.AdminArea) []adminAreaResponse {
	if values == nil {
		return nil
	}
	result := make([]adminAreaResponse, 0, len(values))
	for _, value := range values {
		result = append(result, adminAreaResponse{
			Id:            value.GetId(),
			CountryCode:   value.GetCountryCode(),
			Code:          value.GetCode(),
			Name:          value.GetName(),
			ParentCode:    value.ParentCode,
			Label:         value.GetLabel(),
			Source:        value.GetSource(),
			SourceVersion: value.GetSourceVersion(),
			AdminAreaType: httpdto.EnumString(value.GetAdminAreaType()),
		})
	}
	return result
}

func fromProtoLocalities(values []*locationv1.Locality) []localityResponse {
	if values == nil {
		return nil
	}
	result := make([]localityResponse, 0, len(values))
	for _, value := range values {
		result = append(result, localityResponse{
			Id:            value.GetId(),
			CountryCode:   value.GetCountryCode(),
			AdminAreaCode: value.GetAdminAreaCode(),
			Code:          value.GetCode(),
			Name:          value.GetName(),
			Type:          value.GetType(),
			Label:         value.GetLabel(),
			Source:        value.GetSource(),
			SourceVersion: value.GetSourceVersion(),
		})
	}
	return result
}

func fromProtoSettlements(values []*locationv1.Settlement) []settlementResponse {
	if values == nil {
		return nil
	}
	result := make([]settlementResponse, 0, len(values))
	for _, value := range values {
		result = append(result, settlementResponse{
			Id:            value.GetId(),
			CountryCode:   value.GetCountryCode(),
			AdminAreaCode: value.GetAdminAreaCode(),
			LocalityCode:  value.LocalityCode,
			PostalCode:    value.PostalCode,
			Name:          value.GetName(),
			Type:          value.GetType(),
			Label:         value.GetLabel(),
			Source:        value.GetSource(),
			SourceVersion: value.GetSourceVersion(),
		})
	}
	return result
}

func fromProtoPostalCodeMatches(values []*locationv1.PostalCodeMatch) []postalCodeMatchResponse {
	if values == nil {
		return nil
	}
	result := make([]postalCodeMatchResponse, 0, len(values))
	for _, value := range values {
		result = append(result, postalCodeMatchResponse{
			PostalCode:    value.GetPostalCode(),
			Label:         value.GetLabel(),
			Components:    fromProtoLocationComponents(value.GetComponents()),
			Source:        value.GetSource(),
			SourceVersion: value.GetSourceVersion(),
			Score:         value.GetScore(),
		})
	}
	return result
}

func fromProtoAddressSuggestions(values []*locationv1.AddressSuggestion) []addressSuggestionResponse {
	if values == nil {
		return nil
	}
	result := make([]addressSuggestionResponse, 0, len(values))
	for _, value := range values {
		result = append(result, addressSuggestionResponse{
			Label:         value.GetLabel(),
			Components:    fromProtoLocationComponents(value.GetComponents()),
			Source:        value.GetSource(),
			SourceVersion: value.GetSourceVersion(),
			Score:         value.GetScore(),
		})
	}
	return result
}
