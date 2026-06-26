package application

import (
	"context"

	userv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/user/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/users/ports"
)

type Service struct {
	client ports.Client
}

func NewService(client ports.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Create(ctx context.Context, req *userv1.UserServiceCreateRequest) (*userv1.UserServiceCreateResponse, error) {
	return s.client.Create(ctx, req)
}

func (s *Service) FindByID(ctx context.Context, id string) (*userv1.UserServiceFindByIdResponse, error) {
	return s.client.FindByID(ctx, id)
}

func (s *Service) FindByEmail(ctx context.Context, email string) (*userv1.UserServiceFindByEmailResponse, error) {
	return s.client.FindByEmail(ctx, email)
}

func (s *Service) List(ctx context.Context, req *userv1.UserServiceListRequest) (*userv1.UserServiceListResponse, error) {
	return s.client.List(ctx, req)
}

func (s *Service) UpdatePassword(ctx context.Context, id, newPassword string) (*userv1.UserServiceUpdatePasswordResponse, error) {
	return s.client.UpdatePassword(ctx, id, newPassword)
}

func (s *Service) Delete(ctx context.Context, id string) (*userv1.UserServiceDeleteResponse, error) {
	return s.client.Delete(ctx, id)
}
