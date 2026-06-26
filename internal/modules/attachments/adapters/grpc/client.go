package grpc

import (
	"context"

	filemanagementv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/file_management/v1"
	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/attachments/ports"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
)

type RecordsClient struct {
	client recordv1.AttachmentServiceClient
}

func NewRecords(client recordv1.AttachmentServiceClient) *RecordsClient {
	return &RecordsClient{client: client}
}

func (c *RecordsClient) Create(ctx context.Context, req *recordv1.AttachmentServiceCreateRequest) (*recordv1.AttachmentServiceCreateResponse, error) {
	resp, err := c.client.Create(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *RecordsClient) List(ctx context.Context, req *recordv1.AttachmentServiceListRequest) (*recordv1.AttachmentServiceListResponse, error) {
	resp, err := c.client.List(ctx, req)
	return resp, grpcclient.MapError(err)
}

func (c *RecordsClient) FindByID(ctx context.Context, id string) (*recordv1.AttachmentServiceFindByIdResponse, error) {
	resp, err := c.client.FindById(ctx, &recordv1.AttachmentServiceFindByIdRequest{Id: id})
	return resp, grpcclient.MapError(err)
}

func (c *RecordsClient) Delete(ctx context.Context, id string) (*recordv1.AttachmentServiceDeleteResponse, error) {
	resp, err := c.client.Delete(ctx, &recordv1.AttachmentServiceDeleteRequest{Id: id})
	return resp, grpcclient.MapError(err)
}

type FilesClient struct {
	client filemanagementv1.FileManagementServiceClient
}

func NewFiles(client filemanagementv1.FileManagementServiceClient) *FilesClient {
	return &FilesClient{client: client}
}

func (c *FilesClient) ConfirmUpload(ctx context.Context, req *filemanagementv1.FileManagementServiceConfirmUploadRequest) (*filemanagementv1.FileManagementServiceConfirmUploadResponse, error) {
	resp, err := c.client.ConfirmUpload(ctx, req)
	return resp, grpcclient.MapError(err)
}

var _ ports.RecordsClient = (*RecordsClient)(nil)
var _ ports.FilesClient = (*FilesClient)(nil)
