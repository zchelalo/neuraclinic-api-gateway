package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/catalogs/application"
)

func TestListSexesContract(t *testing.T) {
	handler := NewHandler(application.NewService())
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, passthroughCatalogsMiddleware())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalogs/sexes", nil)
	req.Header.Set("Accept-Language", "es-MX,en;q=0.8")
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	assertCatalogsJSONEqual(t, rec.Body.Bytes(), `{
		"data": [
			{ "value": "SEX_MALE", "label": "Masculino" },
			{ "value": "SEX_FEMALE", "label": "Femenino" },
			{ "value": "SEX_OTHER", "label": "Otro" },
			{ "value": "SEX_PREFER_NOT_TO_SAY", "label": "Prefiero no decirlo" }
		],
		"request_id": ""
	}`)
}

func TestListAppointmentStatusesContract(t *testing.T) {
	handler := NewHandler(application.NewService())
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, passthroughCatalogsMiddleware())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/catalogs/appointment-statuses", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	assertCatalogsJSONEqual(t, rec.Body.Bytes(), `{
		"data": [
			{ "value": "APPOINTMENT_STATUS_SCHEDULED", "label": "Scheduled" },
			{ "value": "APPOINTMENT_STATUS_COMPLETED", "label": "Completed" },
			{ "value": "APPOINTMENT_STATUS_RESCHEDULED", "label": "Rescheduled" },
			{ "value": "APPOINTMENT_STATUS_CANCELLED", "label": "Canceled" },
			{ "value": "APPOINTMENT_STATUS_NO_SHOW", "label": "No show" }
		],
		"request_id": ""
	}`)
}

func passthroughCatalogsMiddleware() middleware.Middleware {
	return func(next http.Handler) http.Handler { return next }
}

func assertCatalogsJSONEqual(t *testing.T, got []byte, expected string) {
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
