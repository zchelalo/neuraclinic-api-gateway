package application

import (
	"context"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/familiograms/ports"
)

type Service struct {
	client ports.Client
}

func NewService(client ports.Client) *Service {
	return &Service{client: client}
}

func (s *Service) FindByPatientID(ctx context.Context, patientID string) (*recordv1.FamiliogramServiceFindByPatientIdResponse, error) {
	return s.client.FindByPatientID(ctx, patientID)
}

func (s *Service) UpdateByPatientID(ctx context.Context, patientID string, req *recordv1.FamiliogramServiceUpdateRequest) (*recordv1.FamiliogramServiceUpdateResponse, error) {
	current, err := s.client.FindByPatientID(ctx, patientID)
	if err != nil {
		return nil, err
	}
	req.Id = current.GetFamiliogram().GetId()
	return s.client.Update(ctx, req)
}
