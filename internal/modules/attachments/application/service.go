package application

import (
	"context"

	filemanagementv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/file_management/v1"
	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/attachments/ports"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
)

type Service struct {
	records ports.RecordsClient
	files   ports.FilesClient
}

func NewService(records ports.RecordsClient, files ports.FilesClient) *Service {
	return &Service{records: records, files: files}
}

func (s *Service) Create(ctx context.Context, req *recordv1.AttachmentServiceCreateRequest) (*recordv1.AttachmentServiceCreateResponse, error) {
	return s.records.Create(ctx, req)
}

func (s *Service) List(ctx context.Context, req *recordv1.AttachmentServiceListRequest) (*recordv1.AttachmentServiceListResponse, error) {
	return s.records.List(ctx, req)
}

func (s *Service) FindByID(ctx context.Context, id string) (*recordv1.AttachmentServiceFindByIdResponse, error) {
	return s.records.FindByID(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id string) (*recordv1.AttachmentServiceDeleteResponse, error) {
	return s.records.Delete(ctx, id)
}

func (s *Service) ConfirmUpload(ctx context.Context, attachmentID string, status sharedv1.FileStatus) (*filemanagementv1.FileManagementServiceConfirmUploadResponse, error) {
	attachment, err := s.records.FindByID(ctx, attachmentID)
	if err != nil {
		return nil, err
	}
	fileID := attachment.GetAttachment().GetFileId()
	if fileID == "" {
		return nil, response.FailedPrecondition("missing_file_id", "attachment has no file id", nil)
	}
	return s.files.ConfirmUpload(ctx, &filemanagementv1.FileManagementServiceConfirmUploadRequest{
		Id:     fileID,
		Status: status,
	})
}
