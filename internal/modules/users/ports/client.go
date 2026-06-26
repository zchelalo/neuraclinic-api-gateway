package ports

import (
	"context"

	userv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/user/v1"
)

type Client interface {
	Create(ctx context.Context, req *userv1.UserServiceCreateRequest) (*userv1.UserServiceCreateResponse, error)
	FindByID(ctx context.Context, id string) (*userv1.UserServiceFindByIdResponse, error)
	FindByEmail(ctx context.Context, email string) (*userv1.UserServiceFindByEmailResponse, error)
	List(ctx context.Context, req *userv1.UserServiceListRequest) (*userv1.UserServiceListResponse, error)
	UpdatePassword(ctx context.Context, id, newPassword string) (*userv1.UserServiceUpdatePasswordResponse, error)
	Delete(ctx context.Context, id string) (*userv1.UserServiceDeleteResponse, error)
}
