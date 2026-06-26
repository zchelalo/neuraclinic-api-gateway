package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	filemanagementv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/file_management/v1"
	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/attachments/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/attachments/ports"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAttachmentContractResponses(t *testing.T) {
	now := timestamppb.New(time.Date(2026, 6, 25, 18, 0, 0, 0, time.UTC))
	records := &fakeAttachmentRecordsClient{
		findByIDResp: &recordv1.AttachmentServiceFindByIdResponse{
			Attachment: &recordv1.Attachment{
				Id:           "att-1",
				FileId:       "file-1",
				MimeType:     "application/pdf",
				DownloadUrl:  strPtr("https://download"),
				ExpiresAt:    now,
				PatientId:    "patient-1",
				CreatedAt:    now,
				UploadStatus: sharedv1.FileStatus_FILE_STATUS_AVAILABLE,
			},
		},
	}
	files := &fakeAttachmentFilesClient{
		confirmResp: &filemanagementv1.FileManagementServiceConfirmUploadResponse{
			Id:          "file-1",
			Status:      sharedv1.FileStatus_FILE_STATUS_AVAILABLE,
			DownloadUrl: strPtr("https://download"),
			ExpiresAt:   now,
		},
	}
	handler := NewHandler(application.NewService(records, files), 10)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, passthroughAttachmentsMiddleware(), func(...string) middleware.Middleware { return passthroughAttachmentsMiddleware() })

	t.Run("find by id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/attachments/att-1", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assertAttachmentsJSONEqual(t, rec.Body.Bytes(), `{
			"data": {
				"id": "att-1",
				"file_id": "file-1",
				"mime_type": "application/pdf",
				"download_url": "https://download",
				"expires_at": "2026-06-25T18:00:00Z",
				"patient_id": "patient-1",
				"created_at": "2026-06-25T18:00:00Z",
				"upload_status": "FILE_STATUS_AVAILABLE"
			},
			"request_id": ""
		}`)
	})

	t.Run("confirm upload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/attachments/att-1/confirm-upload", strings.NewReader(`{"status":"FILE_STATUS_AVAILABLE"}`))
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)

		assertAttachmentsJSONEqual(t, rec.Body.Bytes(), `{
			"data": {
				"id": "file-1",
				"status": "FILE_STATUS_AVAILABLE",
				"download_url": "https://download",
				"expires_at": "2026-06-25T18:00:00Z"
			},
			"request_id": ""
		}`)
	})
}

type fakeAttachmentRecordsClient struct {
	findByIDResp *recordv1.AttachmentServiceFindByIdResponse
}

func (f *fakeAttachmentRecordsClient) Create(context.Context, *recordv1.AttachmentServiceCreateRequest) (*recordv1.AttachmentServiceCreateResponse, error) {
	return nil, nil
}
func (f *fakeAttachmentRecordsClient) List(context.Context, *recordv1.AttachmentServiceListRequest) (*recordv1.AttachmentServiceListResponse, error) {
	return nil, nil
}
func (f *fakeAttachmentRecordsClient) FindByID(context.Context, string) (*recordv1.AttachmentServiceFindByIdResponse, error) {
	return f.findByIDResp, nil
}
func (f *fakeAttachmentRecordsClient) Delete(context.Context, string) (*recordv1.AttachmentServiceDeleteResponse, error) {
	return nil, nil
}

type fakeAttachmentFilesClient struct {
	confirmResp *filemanagementv1.FileManagementServiceConfirmUploadResponse
}

func (f *fakeAttachmentFilesClient) ConfirmUpload(context.Context, *filemanagementv1.FileManagementServiceConfirmUploadRequest) (*filemanagementv1.FileManagementServiceConfirmUploadResponse, error) {
	return f.confirmResp, nil
}

var _ ports.RecordsClient = (*fakeAttachmentRecordsClient)(nil)
var _ ports.FilesClient = (*fakeAttachmentFilesClient)(nil)

func passthroughAttachmentsMiddleware() middleware.Middleware {
	return func(next http.Handler) http.Handler { return next }
}

func strPtr(value string) *string {
	return &value
}

func assertAttachmentsJSONEqual(t *testing.T, got []byte, expected string) {
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
