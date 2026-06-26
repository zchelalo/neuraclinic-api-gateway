package ports

import (
	"context"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
)

type Client interface {
	Create(ctx context.Context, req *recordv1.PatientServiceCreateRequest) (*recordv1.PatientServiceCreateResponse, error)
	List(ctx context.Context, req *recordv1.PatientServiceListRequest) (*recordv1.PatientServiceListResponse, error)
	FindByID(ctx context.Context, id string) (*recordv1.PatientServiceFindByIdResponse, error)
	UpdateIdentification(ctx context.Context, req *recordv1.PatientServiceUpdateIdentificationDataRequest) (*recordv1.PatientServiceUpdateIdentificationDataResponse, error)
	UpdateContact(ctx context.Context, req *recordv1.PatientServiceUpdateContactDetailsRequest) (*recordv1.PatientServiceUpdateContactDetailsResponse, error)
	UpdateAddress(ctx context.Context, req *recordv1.PatientServiceUpdateAddressRequest) (*recordv1.PatientServiceUpdateAddressResponse, error)
}
