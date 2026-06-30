package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/patients/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/patients/ports"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFindByIDContract(t *testing.T) {
	now := timestamppb.New(time.Date(2026, 6, 25, 14, 0, 0, 0, time.UTC))
	handler := NewHandler(application.NewService(&fakePatientsClient{
		findByIDResp: &recordv1.PatientServiceFindByIdResponse{
			Patient: &recordv1.Patient{
				Id:            "patient-1",
				FirstName:     "Ana",
				FirstLastName: "Lopez",
				BirthDate:     &date.Date{Year: 1994, Month: 5, Day: 10},
				BirthCountry:  "MX",
				BirthProvince: "JAL",
				BirthCity:     "Guadalajara",
				Sex:           recordv1.Sex_SEX_FEMALE,
				MaritalStatus: recordv1.MaritalStatus_MARITAL_STATUS_SINGLE,
				Phone:         "5551234567",
				Email:         "ana@example.com",
				Address: &recordv1.Address{
					Id:           "addr-1",
					Country:      "MX",
					Province:     "JAL",
					City:         "Guadalajara",
					PostalCode:   "44100",
					Neighborhood: "Centro",
					Street:       "Juarez",
					StreetNumber: "10",
					CreatedAt:    now,
				},
				PsychologistId: "psy-1",
				CreatedAt:      now,
			},
		},
	}), 10)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, passthroughPatientsMiddleware(), func(...string) middleware.Middleware { return passthroughPatientsMiddleware() })

	req := httptest.NewRequest(http.MethodGet, "/api/v1/patients/patient-1", nil)
	req.Header.Set("Accept-Language", "es-MX,en;q=0.8")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	assertPatientsJSONEqual(t, rec.Body.Bytes(), `{
		"data": {
			"id": "patient-1",
			"first_name": "Ana",
			"first_last_name": "Lopez",
			"birth_date": {"year": 1994, "month": 5, "day": 10},
			"birth_country": "MX",
			"birth_province": "JAL",
			"birth_city": "Guadalajara",
			"sex": {
				"value": "SEX_FEMALE",
				"label": "Femenino"
			},
			"marital_status": {
				"value": "MARITAL_STATUS_SINGLE",
				"label": "Soltero"
			},
			"phone": "5551234567",
			"email": "ana@example.com",
			"address": {
				"id": "addr-1",
				"country": "MX",
				"province": "JAL",
				"city": "Guadalajara",
				"postal_code": "44100",
				"neighborhood": "Centro",
				"street": "Juarez",
				"street_number": "10",
				"created_at": "2026-06-25T14:00:00Z"
			},
			"psychologist_id": "psy-1",
			"created_at": "2026-06-25T14:00:00Z"
		},
		"request_id": ""
	}`)
}

type fakePatientsClient struct {
	findByIDResp *recordv1.PatientServiceFindByIdResponse
}

func (f *fakePatientsClient) Create(context.Context, *recordv1.PatientServiceCreateRequest) (*recordv1.PatientServiceCreateResponse, error) {
	return nil, nil
}
func (f *fakePatientsClient) List(context.Context, *recordv1.PatientServiceListRequest) (*recordv1.PatientServiceListResponse, error) {
	return nil, nil
}
func (f *fakePatientsClient) FindByID(context.Context, string) (*recordv1.PatientServiceFindByIdResponse, error) {
	return f.findByIDResp, nil
}
func (f *fakePatientsClient) UpdateIdentification(context.Context, *recordv1.PatientServiceUpdateIdentificationDataRequest) (*recordv1.PatientServiceUpdateIdentificationDataResponse, error) {
	return nil, nil
}
func (f *fakePatientsClient) UpdateContact(context.Context, *recordv1.PatientServiceUpdateContactDetailsRequest) (*recordv1.PatientServiceUpdateContactDetailsResponse, error) {
	return nil, nil
}
func (f *fakePatientsClient) UpdateAddress(context.Context, *recordv1.PatientServiceUpdateAddressRequest) (*recordv1.PatientServiceUpdateAddressResponse, error) {
	return nil, nil
}

var _ ports.Client = (*fakePatientsClient)(nil)

func passthroughPatientsMiddleware() middleware.Middleware {
	return func(next http.Handler) http.Handler { return next }
}

func assertPatientsJSONEqual(t *testing.T, got []byte, expected string) {
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
