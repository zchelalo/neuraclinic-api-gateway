package application

import (
	"context"
	"testing"

	filemanagementv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/file_management/v1"
	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
)

func TestConfirmUploadResolvesAttachmentBeforeCallingFiles(t *testing.T) {
	records := &fakeRecords{
		findResp: &recordv1.AttachmentServiceFindByIdResponse{
			Attachment: &recordv1.Attachment{
				Id:     "attachment-1",
				FileId: "file-1",
			},
		},
	}
	files := &fakeFiles{
		resp: &filemanagementv1.FileManagementServiceConfirmUploadResponse{
			Id:     "file-1",
			Status: sharedv1.FileStatus_FILE_STATUS_AVAILABLE,
		},
	}
	service := NewService(records, files)

	resp, err := service.ConfirmUpload(context.Background(), "attachment-1", sharedv1.FileStatus_FILE_STATUS_AVAILABLE)
	if err != nil {
		t.Fatalf("ConfirmUpload returned error: %v", err)
	}
	if resp.GetId() != "file-1" {
		t.Fatalf("expected file id file-1, got %q", resp.GetId())
	}
	if files.lastReq == nil || files.lastReq.GetId() != "file-1" {
		t.Fatal("expected confirm upload to call file-management with resolved file id")
	}
}

type fakeRecords struct {
	findResp *recordv1.AttachmentServiceFindByIdResponse
}

func (f *fakeRecords) Create(context.Context, *recordv1.AttachmentServiceCreateRequest) (*recordv1.AttachmentServiceCreateResponse, error) {
	return nil, nil
}

func (f *fakeRecords) List(context.Context, *recordv1.AttachmentServiceListRequest) (*recordv1.AttachmentServiceListResponse, error) {
	return nil, nil
}

func (f *fakeRecords) FindByID(context.Context, string) (*recordv1.AttachmentServiceFindByIdResponse, error) {
	return f.findResp, nil
}

func (f *fakeRecords) Delete(context.Context, string) (*recordv1.AttachmentServiceDeleteResponse, error) {
	return nil, nil
}

type fakeFiles struct {
	resp    *filemanagementv1.FileManagementServiceConfirmUploadResponse
	lastReq *filemanagementv1.FileManagementServiceConfirmUploadRequest
}

func (f *fakeFiles) ConfirmUpload(_ context.Context, req *filemanagementv1.FileManagementServiceConfirmUploadRequest) (*filemanagementv1.FileManagementServiceConfirmUploadResponse, error) {
	f.lastReq = req
	return f.resp, nil
}
