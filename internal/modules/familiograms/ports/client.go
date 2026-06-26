package ports

import (
	"context"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
)

type Client interface {
	FindByPatientID(ctx context.Context, patientID string) (*recordv1.FamiliogramServiceFindByPatientIdResponse, error)
	Update(ctx context.Context, req *recordv1.FamiliogramServiceUpdateRequest) (*recordv1.FamiliogramServiceUpdateResponse, error)
}
