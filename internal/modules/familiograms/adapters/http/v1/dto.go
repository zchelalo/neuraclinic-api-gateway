package v1

import (
	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/httpdto"
)

type familiogramResponse struct {
	Id        string         `json:"id,omitempty"`
	Data      map[string]any `json:"data,omitempty"`
	PatientID string         `json:"patient_id,omitempty"`
	CreatedAt *string        `json:"created_at,omitempty"`
	UpdatedAt *string        `json:"updated_at,omitempty"`
}

func fromProtoFamiliogram(value *recordv1.Familiogram) *familiogramResponse {
	if value == nil {
		return nil
	}
	return &familiogramResponse{
		Id:        value.GetId(),
		Data:      httpdto.Struct(value.GetData()),
		PatientID: value.GetPatientId(),
		CreatedAt: httpdto.Timestamp(value.GetCreatedAt()),
		UpdatedAt: httpdto.Timestamp(value.GetUpdatedAt()),
	}
}
