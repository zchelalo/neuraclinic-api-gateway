package ports

import (
	"context"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
)

type Client interface {
	Create(ctx context.Context, req *recordv1.NoteServiceCreateRequest) (*recordv1.NoteServiceCreateResponse, error)
	List(ctx context.Context, req *recordv1.NoteServiceListRequest) (*recordv1.NoteServiceListResponse, error)
	FindByID(ctx context.Context, id string) (*recordv1.NoteServiceFindByIdResponse, error)
	Update(ctx context.Context, req *recordv1.NoteServiceUpdateRequest) (*recordv1.NoteServiceUpdateResponse, error)
	Delete(ctx context.Context, id string) (*recordv1.NoteServiceDeleteResponse, error)
}
