package application

import (
	"context"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/appointments/ports"
)

type Service struct {
	client ports.Client
}

func NewService(client ports.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Create(ctx context.Context, req *recordv1.AppointmentServiceCreateRequest) (*recordv1.AppointmentServiceCreateResponse, error) {
	return s.client.Create(ctx, req)
}

func (s *Service) List(ctx context.Context, req *recordv1.AppointmentServiceListRequest) (*recordv1.AppointmentServiceListResponse, error) {
	return s.client.List(ctx, req)
}

func (s *Service) FindByID(ctx context.Context, id string) (*recordv1.AppointmentServiceFindByIdResponse, error) {
	return s.client.FindByID(ctx, id)
}

func (s *Service) Reschedule(ctx context.Context, req *recordv1.AppointmentServiceRescheduleRequest) (*recordv1.AppointmentServiceRescheduleResponse, error) {
	return s.client.Reschedule(ctx, req)
}

func (s *Service) UpdateStatus(ctx context.Context, req *recordv1.AppointmentServiceUpdateStatusRequest) (*recordv1.AppointmentServiceUpdateStatusResponse, error) {
	return s.client.UpdateStatus(ctx, req)
}
