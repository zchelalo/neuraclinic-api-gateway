package v1

import (
	"net/http"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/patients/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/server/httpx"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/language"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/meta"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/parse"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
	"google.golang.org/genproto/googleapis/type/date"
)

const (
	permPatientCreate = "PERMISSION_KEY_PATIENT_CREATE"
	permPatientView   = "PERMISSION_KEY_PATIENT_VIEW"
	permPatientEdit   = "PERMISSION_KEY_PATIENT_EDIT"
)

type Handler struct {
	service      *application.Service
	defaultLimit int32
}

func NewHandler(service *application.Service, defaultLimit int32) *Handler {
	return &Handler{service: service, defaultLimit: defaultLimit}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, auth middleware.Middleware, permissions func(...string) middleware.Middleware) {
	mux.Handle("POST /api/v1/patients", middleware.Chain(http.HandlerFunc(h.create), auth, permissions(permPatientCreate)))
	mux.Handle("GET /api/v1/patients", middleware.Chain(http.HandlerFunc(h.list), auth, permissions(permPatientView)))
	mux.Handle("GET /api/v1/patients/{id}", middleware.Chain(http.HandlerFunc(h.findByID), auth, permissions(permPatientView)))
	mux.Handle("PATCH /api/v1/patients/{id}/identification", middleware.Chain(http.HandlerFunc(h.updateIdentification), auth, permissions(permPatientEdit)))
	mux.Handle("PATCH /api/v1/patients/{id}/contact", middleware.Chain(http.HandlerFunc(h.updateContact), auth, permissions(permPatientEdit)))
	mux.Handle("PATCH /api/v1/patients/{id}/address", middleware.Chain(http.HandlerFunc(h.updateAddress), auth, permissions(permPatientEdit)))
}

type createPatientRequest struct {
	FirstName      string `json:"first_name"`
	MiddleName     string `json:"middle_name"`
	FirstLastName  string `json:"first_last_name"`
	SecondLastName string `json:"second_last_name"`
	BirthDate      string `json:"birth_date"`
	BirthCountry   string `json:"birth_country"`
	BirthProvince  string `json:"birth_province"`
	BirthCity      string `json:"birth_city"`
	Sex            string `json:"sex"`
	MaritalStatus  string `json:"marital_status"`
	Occupation     string `json:"occupation"`
	Religion       string `json:"religion"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Country        string `json:"country"`
	Province       string `json:"province"`
	City           string `json:"city"`
	PostalCode     string `json:"postal_code"`
	Neighborhood   string `json:"neighborhood"`
	Street         string `json:"street"`
	StreetNumber   string `json:"street_number"`
	UnitNumber     string `json:"unit_number"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var body createPatientRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	birthDate, err := parse.Date(body.BirthDate)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	sex, err := parse.Sex(body.Sex)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	maritalStatus, err := parse.MaritalStatus(body.MaritalStatus)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.Create(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &recordv1.PatientServiceCreateRequest{
		FirstName:      body.FirstName,
		MiddleName:     parse.OptionalString(body.MiddleName),
		FirstLastName:  body.FirstLastName,
		SecondLastName: parse.OptionalString(body.SecondLastName),
		BirthDate:      birthDate,
		BirthCountry:   body.BirthCountry,
		BirthProvince:  body.BirthProvince,
		BirthCity:      body.BirthCity,
		Sex:            sex,
		MaritalStatus:  maritalStatus,
		Occupation:     parse.OptionalString(body.Occupation),
		Religion:       parse.OptionalString(body.Religion),
		Phone:          body.Phone,
		Email:          body.Email,
		Country:        body.Country,
		Province:       body.Province,
		City:           body.City,
		PostalCode:     body.PostalCode,
		Neighborhood:   body.Neighborhood,
		Street:         body.Street,
		StreetNumber:   body.StreetNumber,
		UnitNumber:     parse.OptionalString(body.UnitNumber),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusCreated, fromProtoPatientSummary(resp.GetPatient()), nil)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	pagination, err := parse.CursorPagination(r, h.defaultLimit)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	withPending, err := parse.QueryOptionalBool(r, "with_pending_appointments")
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	withNoAppointments, err := parse.QueryOptionalBool(r, "with_no_appointments")
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	everHad, err := parse.QueryOptionalBool(r, "ever_had_appointments")
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.List(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &recordv1.PatientServiceListRequest{
		Pagination:              pagination,
		WithPendingAppointments: withPending,
		WithNoAppointments:      withNoAppointments,
		EverHadAppointments:     everHad,
		SearchQuery:             parse.OptionalString(parse.QueryString(r, "search_query")),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoPatientSummaries(resp.GetPatients()), meta.FromProto(resp.GetMeta()))
}

func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.FindByID(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), r.PathValue("id"))
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoPatient(resp.GetPatient(), language.ResolveRequest(r)), nil)
}

type updateIdentificationRequest struct {
	FirstName      *string `json:"first_name,omitempty"`
	MiddleName     *string `json:"middle_name,omitempty"`
	FirstLastName  *string `json:"first_last_name,omitempty"`
	SecondLastName *string `json:"second_last_name,omitempty"`
	BirthDate      *string `json:"birth_date,omitempty"`
	Sex            *string `json:"sex,omitempty"`
	BirthCountry   *string `json:"birth_country,omitempty"`
	BirthProvince  *string `json:"birth_province,omitempty"`
	BirthCity      *string `json:"birth_city,omitempty"`
	Occupation     *string `json:"occupation,omitempty"`
	MaritalStatus  *string `json:"marital_status,omitempty"`
	Religion       *string `json:"religion,omitempty"`
}

func (h *Handler) updateIdentification(w http.ResponseWriter, r *http.Request) {
	var body updateIdentificationRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	var birthDate *date.Date
	var err error
	if body.BirthDate != nil {
		birthDate, err = parse.Date(*body.BirthDate)
		if err != nil {
			response.WriteError(w, r, err)
			return
		}
	}
	var sex *recordv1.Sex
	if body.Sex != nil {
		value, err := parse.Sex(*body.Sex)
		if err != nil {
			response.WriteError(w, r, err)
			return
		}
		sex = &value
	}
	var maritalStatus *recordv1.MaritalStatus
	if body.MaritalStatus != nil {
		value, err := parse.MaritalStatus(*body.MaritalStatus)
		if err != nil {
			response.WriteError(w, r, err)
			return
		}
		maritalStatus = &value
	}
	resp, err := h.service.UpdateIdentification(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &recordv1.PatientServiceUpdateIdentificationDataRequest{
		Id:             r.PathValue("id"),
		FirstName:      body.FirstName,
		MiddleName:     body.MiddleName,
		FirstLastName:  body.FirstLastName,
		SecondLastName: body.SecondLastName,
		BirthDate:      birthDate,
		Sex:            sex,
		BirthCountry:   body.BirthCountry,
		BirthProvince:  body.BirthProvince,
		BirthCity:      body.BirthCity,
		Occupation:     body.Occupation,
		MaritalStatus:  maritalStatus,
		Religion:       body.Religion,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoPatient(resp.GetPatient(), language.ResolveRequest(r)), nil)
}

type updateContactRequest struct {
	Phone *string `json:"phone,omitempty"`
	Email *string `json:"email,omitempty"`
}

func (h *Handler) updateContact(w http.ResponseWriter, r *http.Request) {
	var body updateContactRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.UpdateContact(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &recordv1.PatientServiceUpdateContactDetailsRequest{
		Id:    r.PathValue("id"),
		Phone: body.Phone,
		Email: body.Email,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoPatient(resp.GetPatient(), language.ResolveRequest(r)), nil)
}

type updateAddressRequest struct {
	Country      *string `json:"country,omitempty"`
	Province     *string `json:"province,omitempty"`
	City         *string `json:"city,omitempty"`
	PostalCode   *string `json:"postal_code,omitempty"`
	Neighborhood *string `json:"neighborhood,omitempty"`
	Street       *string `json:"street,omitempty"`
	StreetNumber *string `json:"street_number,omitempty"`
	UnitNumber   *string `json:"unit_number,omitempty"`
}

func (h *Handler) updateAddress(w http.ResponseWriter, r *http.Request) {
	var body updateAddressRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.UpdateAddress(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &recordv1.PatientServiceUpdateAddressRequest{
		Id:           r.PathValue("id"),
		Country:      body.Country,
		Province:     body.Province,
		City:         body.City,
		PostalCode:   body.PostalCode,
		Neighborhood: body.Neighborhood,
		Street:       body.Street,
		StreetNumber: body.StreetNumber,
		UnitNumber:   body.UnitNumber,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoPatient(resp.GetPatient(), language.ResolveRequest(r)), nil)
}
