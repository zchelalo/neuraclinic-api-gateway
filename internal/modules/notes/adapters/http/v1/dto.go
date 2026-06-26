package v1

import (
	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/httpdto"
)

type noteResponse struct {
	Id            string  `json:"id,omitempty"`
	PatientID     string  `json:"patient_id,omitempty"`
	AppointmentID *string `json:"appointment_id,omitempty"`
	Title         *string `json:"title,omitempty"`
	ContentHTML   string  `json:"content_html,omitempty"`
	ContentText   string  `json:"content_text,omitempty"`
	CreatedAt     *string `json:"created_at,omitempty"`
	UpdatedAt     *string `json:"updated_at,omitempty"`
	DeletedAt     *string `json:"deleted_at,omitempty"`
}

type noteSummaryResponse struct {
	Id            string  `json:"id,omitempty"`
	PatientID     string  `json:"patient_id,omitempty"`
	AppointmentID *string `json:"appointment_id,omitempty"`
	Title         *string `json:"title,omitempty"`
	CreatedAt     *string `json:"created_at,omitempty"`
	UpdatedAt     *string `json:"updated_at,omitempty"`
	DeletedAt     *string `json:"deleted_at,omitempty"`
}

func fromProtoNote(value *recordv1.Note) *noteResponse {
	if value == nil {
		return nil
	}
	return &noteResponse{
		Id:            value.GetId(),
		PatientID:     value.GetPatientId(),
		AppointmentID: value.AppointmentId,
		Title:         value.Title,
		ContentHTML:   value.GetContentHtml(),
		ContentText:   value.GetContentText(),
		CreatedAt:     httpdto.Timestamp(value.GetCreatedAt()),
		UpdatedAt:     httpdto.Timestamp(value.GetUpdatedAt()),
		DeletedAt:     httpdto.Timestamp(value.GetDeletedAt()),
	}
}

func fromProtoNoteSummaries(values []*recordv1.NoteSummary) []noteSummaryResponse {
	if values == nil {
		return nil
	}
	result := make([]noteSummaryResponse, 0, len(values))
	for _, value := range values {
		result = append(result, noteSummaryResponse{
			Id:            value.GetId(),
			PatientID:     value.GetPatientId(),
			AppointmentID: value.AppointmentId,
			Title:         value.Title,
			CreatedAt:     httpdto.Timestamp(value.GetCreatedAt()),
			UpdatedAt:     httpdto.Timestamp(value.GetUpdatedAt()),
			DeletedAt:     httpdto.Timestamp(value.GetDeletedAt()),
		})
	}
	return result
}
