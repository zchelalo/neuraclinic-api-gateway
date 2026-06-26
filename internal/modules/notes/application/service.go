package application

import (
	"context"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/notes/ports"
)

type Service struct {
	client ports.Client
}

func NewService(client ports.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Create(ctx context.Context, req *recordv1.NoteServiceCreateRequest) (*recordv1.NoteServiceCreateResponse, error) {
	return s.client.Create(ctx, req)
}

func (s *Service) List(ctx context.Context, req *recordv1.NoteServiceListRequest) (*recordv1.NoteServiceListResponse, error) {
	return s.client.List(ctx, req)
}

func (s *Service) FindByID(ctx context.Context, id string) (*recordv1.NoteServiceFindByIdResponse, error) {
	return s.client.FindByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, req *recordv1.NoteServiceUpdateRequest) (*recordv1.NoteServiceUpdateResponse, error) {
	return s.client.Update(ctx, req)
}

func (s *Service) Delete(ctx context.Context, id string) (*recordv1.NoteServiceDeleteResponse, error) {
	return s.client.Delete(ctx, id)
}
