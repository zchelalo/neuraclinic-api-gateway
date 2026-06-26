package application

import (
	"context"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/patients/ports"
)

type Service struct {
	client ports.Client
}

func NewService(client ports.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Create(ctx context.Context, req *recordv1.PatientServiceCreateRequest) (*recordv1.PatientServiceCreateResponse, error) {
	return s.client.Create(ctx, req)
}

func (s *Service) List(ctx context.Context, req *recordv1.PatientServiceListRequest) (*recordv1.PatientServiceListResponse, error) {
	return s.client.List(ctx, req)
}

func (s *Service) FindByID(ctx context.Context, id string) (*recordv1.PatientServiceFindByIdResponse, error) {
	return s.client.FindByID(ctx, id)
}

func (s *Service) UpdateIdentification(ctx context.Context, req *recordv1.PatientServiceUpdateIdentificationDataRequest) (*recordv1.PatientServiceUpdateIdentificationDataResponse, error) {
	return s.client.UpdateIdentification(ctx, req)
}

func (s *Service) UpdateContact(ctx context.Context, req *recordv1.PatientServiceUpdateContactDetailsRequest) (*recordv1.PatientServiceUpdateContactDetailsResponse, error) {
	return s.client.UpdateContact(ctx, req)
}

func (s *Service) UpdateAddress(ctx context.Context, req *recordv1.PatientServiceUpdateAddressRequest) (*recordv1.PatientServiceUpdateAddressResponse, error) {
	return s.client.UpdateAddress(ctx, req)
}
