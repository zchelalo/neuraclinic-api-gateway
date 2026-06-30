package v1

import (
	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	catalogsapplication "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/catalogs/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/httpdto"
)

type addressResponse struct {
	Id           string  `json:"id,omitempty"`
	Country      string  `json:"country,omitempty"`
	Province     string  `json:"province,omitempty"`
	City         string  `json:"city,omitempty"`
	PostalCode   string  `json:"postal_code,omitempty"`
	Neighborhood string  `json:"neighborhood,omitempty"`
	Street       string  `json:"street,omitempty"`
	StreetNumber string  `json:"street_number,omitempty"`
	UnitNumber   *string `json:"unit_number,omitempty"`
	CreatedAt    *string `json:"created_at,omitempty"`
	UpdatedAt    *string `json:"updated_at,omitempty"`
	DeletedAt    *string `json:"deleted_at,omitempty"`
}

type enumValueResponse struct {
	Value string `json:"value,omitempty"`
	Label string `json:"label,omitempty"`
}

type patientResponse struct {
	Id             string            `json:"id,omitempty"`
	FirstName      string            `json:"first_name,omitempty"`
	MiddleName     *string           `json:"middle_name,omitempty"`
	FirstLastName  string            `json:"first_last_name,omitempty"`
	SecondLastName *string           `json:"second_last_name,omitempty"`
	BirthDate      *httpdto.Date     `json:"birth_date,omitempty"`
	BirthCountry   string            `json:"birth_country,omitempty"`
	BirthProvince  string            `json:"birth_province,omitempty"`
	BirthCity      string            `json:"birth_city,omitempty"`
	Sex            enumValueResponse `json:"sex,omitempty"`
	MaritalStatus  enumValueResponse `json:"marital_status,omitempty"`
	Occupation     *string           `json:"occupation,omitempty"`
	Religion       *string           `json:"religion,omitempty"`
	Phone          string            `json:"phone,omitempty"`
	Email          string            `json:"email,omitempty"`
	Address        *addressResponse  `json:"address,omitempty"`
	PsychologistID string            `json:"psychologist_id,omitempty"`
	CreatedAt      *string           `json:"created_at,omitempty"`
	UpdatedAt      *string           `json:"updated_at,omitempty"`
	DeletedAt      *string           `json:"deleted_at,omitempty"`
}

type patientSummaryResponse struct {
	Id             string  `json:"id,omitempty"`
	FirstName      string  `json:"first_name,omitempty"`
	MiddleName     *string `json:"middle_name,omitempty"`
	FirstLastName  string  `json:"first_last_name,omitempty"`
	SecondLastName *string `json:"second_last_name,omitempty"`
	BirthDate      string  `json:"birth_date,omitempty"`
	Email          string  `json:"email,omitempty"`
	Phone          string  `json:"phone,omitempty"`
}

func fromProtoPatientSummary(value *recordv1.PatientSummary) *patientSummaryResponse {
	if value == nil {
		return nil
	}
	return &patientSummaryResponse{
		Id:             value.GetId(),
		FirstName:      value.GetFirstName(),
		MiddleName:     value.MiddleName,
		FirstLastName:  value.GetFirstLastName(),
		SecondLastName: value.SecondLastName,
		BirthDate:      value.GetBirthDate(),
		Email:          value.GetEmail(),
		Phone:          value.GetPhone(),
	}
}

func fromProtoAddress(value *recordv1.Address) *addressResponse {
	if value == nil {
		return nil
	}
	return &addressResponse{
		Id:           value.GetId(),
		Country:      value.GetCountry(),
		Province:     value.GetProvince(),
		City:         value.GetCity(),
		PostalCode:   value.GetPostalCode(),
		Neighborhood: value.GetNeighborhood(),
		Street:       value.GetStreet(),
		StreetNumber: value.GetStreetNumber(),
		UnitNumber:   value.UnitNumber,
		CreatedAt:    httpdto.Timestamp(value.GetCreatedAt()),
		UpdatedAt:    httpdto.Timestamp(value.GetUpdatedAt()),
		DeletedAt:    httpdto.Timestamp(value.GetDeletedAt()),
	}
}

func fromProtoPatient(value *recordv1.Patient, locale string) *patientResponse {
	if value == nil {
		return nil
	}
	return &patientResponse{
		Id:             value.GetId(),
		FirstName:      value.GetFirstName(),
		MiddleName:     value.MiddleName,
		FirstLastName:  value.GetFirstLastName(),
		SecondLastName: value.SecondLastName,
		BirthDate:      httpdto.DateValue(value.GetBirthDate()),
		BirthCountry:   value.GetBirthCountry(),
		BirthProvince:  value.GetBirthProvince(),
		BirthCity:      value.GetBirthCity(),
		Sex: enumValueResponse{
			Value: httpdto.EnumString(value.GetSex()),
			Label: catalogsapplication.LocalizedSexLabel(locale, value.GetSex()),
		},
		MaritalStatus: enumValueResponse{
			Value: httpdto.EnumString(value.GetMaritalStatus()),
			Label: catalogsapplication.LocalizedMaritalStatusLabel(
				locale,
				value.GetMaritalStatus(),
			),
		},
		Occupation:     value.Occupation,
		Religion:       value.Religion,
		Phone:          value.GetPhone(),
		Email:          value.GetEmail(),
		Address:        fromProtoAddress(value.GetAddress()),
		PsychologistID: value.GetPsychologistId(),
		CreatedAt:      httpdto.Timestamp(value.GetCreatedAt()),
		UpdatedAt:      httpdto.Timestamp(value.GetUpdatedAt()),
		DeletedAt:      httpdto.Timestamp(value.GetDeletedAt()),
	}
}

func fromProtoPatientSummaries(values []*recordv1.PatientSummary) []patientSummaryResponse {
	if values == nil {
		return nil
	}
	result := make([]patientSummaryResponse, 0, len(values))
	for _, value := range values {
		mapped := fromProtoPatientSummary(value)
		if mapped == nil {
			result = append(result, patientSummaryResponse{})
			continue
		}
		result = append(result, *mapped)
	}
	return result
}
