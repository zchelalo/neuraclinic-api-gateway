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
	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/appointments/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/appointments/ports"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFindByIDContract(t *testing.T) {
	start := timestamppb.New(time.Date(2026, 6, 25, 15, 0, 0, 0, time.UTC))
	end := timestamppb.New(time.Date(2026, 6, 25, 16, 0, 0, 0, time.UTC))
	handler := NewHandler(application.NewService(&fakeAppointmentsClient{
		findByIDResp: &recordv1.AppointmentServiceFindByIdResponse{
			Appointment: &recordv1.Appointment{
				Id:                           "appt-1",
				StartTime:                    start,
				EndTime:                      end,
				Reason:                       "Therapy",
				Status:                       sharedv1.AppointmentStatus_APPOINTMENT_STATUS_SCHEDULED,
				PatientId:                    "patient-1",
				RescheduledFromAppointmentId: strPtr("appt-0"),
				CreatedAt:                    start,
			},
		},
	}), 10)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, passthroughAppointmentsMiddleware(), func(...string) middleware.Middleware { return passthroughAppointmentsMiddleware() })

	req := httptest.NewRequest(http.MethodGet, "/api/v1/appointments/appt-1", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	assertAppointmentsJSONEqual(t, rec.Body.Bytes(), `{
		"data": {
			"id": "appt-1",
			"start_time": "2026-06-25T15:00:00Z",
			"end_time": "2026-06-25T16:00:00Z",
			"reason": "Therapy",
			"status": "APPOINTMENT_STATUS_SCHEDULED",
			"patient_id": "patient-1",
			"rescheduled_from_appointment_id": "appt-0",
			"created_at": "2026-06-25T15:00:00Z"
		},
		"request_id": ""
	}`)
}

type fakeAppointmentsClient struct {
	findByIDResp *recordv1.AppointmentServiceFindByIdResponse
}

func (f *fakeAppointmentsClient) Create(context.Context, *recordv1.AppointmentServiceCreateRequest) (*recordv1.AppointmentServiceCreateResponse, error) {
	return nil, nil
}
func (f *fakeAppointmentsClient) List(context.Context, *recordv1.AppointmentServiceListRequest) (*recordv1.AppointmentServiceListResponse, error) {
	return nil, nil
}
func (f *fakeAppointmentsClient) FindByID(context.Context, string) (*recordv1.AppointmentServiceFindByIdResponse, error) {
	return f.findByIDResp, nil
}
func (f *fakeAppointmentsClient) Reschedule(context.Context, *recordv1.AppointmentServiceRescheduleRequest) (*recordv1.AppointmentServiceRescheduleResponse, error) {
	return nil, nil
}
func (f *fakeAppointmentsClient) UpdateStatus(context.Context, *recordv1.AppointmentServiceUpdateStatusRequest) (*recordv1.AppointmentServiceUpdateStatusResponse, error) {
	return nil, nil
}

var _ ports.Client = (*fakeAppointmentsClient)(nil)

func passthroughAppointmentsMiddleware() middleware.Middleware {
	return func(next http.Handler) http.Handler { return next }
}

func strPtr(value string) *string {
	return &value
}

func assertAppointmentsJSONEqual(t *testing.T, got []byte, expected string) {
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
