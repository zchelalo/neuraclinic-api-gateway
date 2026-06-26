package ports

import (
	"context"

	filemanagementv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/file_management/v1"
	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
)

type RecordsClient interface {
	Create(ctx context.Context, req *recordv1.AttachmentServiceCreateRequest) (*recordv1.AttachmentServiceCreateResponse, error)
	List(ctx context.Context, req *recordv1.AttachmentServiceListRequest) (*recordv1.AttachmentServiceListResponse, error)
	FindByID(ctx context.Context, id string) (*recordv1.AttachmentServiceFindByIdResponse, error)
	Delete(ctx context.Context, id string) (*recordv1.AttachmentServiceDeleteResponse, error)
}

type FilesClient interface {
	ConfirmUpload(ctx context.Context, req *filemanagementv1.FileManagementServiceConfirmUploadRequest) (*filemanagementv1.FileManagementServiceConfirmUploadResponse, error)
}
