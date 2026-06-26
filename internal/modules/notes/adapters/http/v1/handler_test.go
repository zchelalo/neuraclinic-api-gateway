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
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/notes/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/notes/ports"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestFindByIDContract(t *testing.T) {
	now := timestamppb.New(time.Date(2026, 6, 25, 17, 0, 0, 0, time.UTC))
	handler := NewHandler(application.NewService(&fakeNotesClient{
		findByIDResp: &recordv1.NoteServiceFindByIdResponse{
			Note: &recordv1.Note{
				Id:            "note-1",
				PatientId:     "patient-1",
				AppointmentId: strPtr("appt-1"),
				Title:         strPtr("Progress note"),
				ContentHtml:   "<p>Hello</p>",
				ContentText:   "Hello",
				CreatedAt:     now,
			},
		},
	}), 10)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, passthroughNotesMiddleware(), func(...string) middleware.Middleware { return passthroughNotesMiddleware() })

	req := httptest.NewRequest(http.MethodGet, "/api/v1/notes/note-1", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	assertNotesJSONEqual(t, rec.Body.Bytes(), `{
		"data": {
			"id": "note-1",
			"patient_id": "patient-1",
			"appointment_id": "appt-1",
			"title": "Progress note",
			"content_html": "<p>Hello</p>",
			"content_text": "Hello",
			"created_at": "2026-06-25T17:00:00Z"
		},
		"request_id": ""
	}`)
}

type fakeNotesClient struct {
	findByIDResp *recordv1.NoteServiceFindByIdResponse
}

func (f *fakeNotesClient) Create(context.Context, *recordv1.NoteServiceCreateRequest) (*recordv1.NoteServiceCreateResponse, error) {
	return nil, nil
}
func (f *fakeNotesClient) List(context.Context, *recordv1.NoteServiceListRequest) (*recordv1.NoteServiceListResponse, error) {
	return nil, nil
}
func (f *fakeNotesClient) FindByID(context.Context, string) (*recordv1.NoteServiceFindByIdResponse, error) {
	return f.findByIDResp, nil
}
func (f *fakeNotesClient) Update(context.Context, *recordv1.NoteServiceUpdateRequest) (*recordv1.NoteServiceUpdateResponse, error) {
	return nil, nil
}
func (f *fakeNotesClient) Delete(context.Context, string) (*recordv1.NoteServiceDeleteResponse, error) {
	return nil, nil
}

var _ ports.Client = (*fakeNotesClient)(nil)

func passthroughNotesMiddleware() middleware.Middleware {
	return func(next http.Handler) http.Handler { return next }
}

func strPtr(value string) *string {
	return &value
}

func assertNotesJSONEqual(t *testing.T, got []byte, expected string) {
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
