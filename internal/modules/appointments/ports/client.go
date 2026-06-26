package ports

import (
	"context"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
)

type Client interface {
	Create(ctx context.Context, req *recordv1.AppointmentServiceCreateRequest) (*recordv1.AppointmentServiceCreateResponse, error)
	List(ctx context.Context, req *recordv1.AppointmentServiceListRequest) (*recordv1.AppointmentServiceListResponse, error)
	FindByID(ctx context.Context, id string) (*recordv1.AppointmentServiceFindByIdResponse, error)
	Reschedule(ctx context.Context, req *recordv1.AppointmentServiceRescheduleRequest) (*recordv1.AppointmentServiceRescheduleResponse, error)
	UpdateStatus(ctx context.Context, req *recordv1.AppointmentServiceUpdateStatusRequest) (*recordv1.AppointmentServiceUpdateStatusResponse, error)
}
