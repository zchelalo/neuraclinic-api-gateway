package v1

import (
	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/httpdto"
)

type appointmentResponse struct {
	Id                           string  `json:"id,omitempty"`
	StartTime                    *string `json:"start_time,omitempty"`
	EndTime                      *string `json:"end_time,omitempty"`
	Reason                       string  `json:"reason,omitempty"`
	Status                       string  `json:"status,omitempty"`
	PatientID                    string  `json:"patient_id,omitempty"`
	CancelledByUserID            *string `json:"cancelled_by_user_id,omitempty"`
	RescheduledFromAppointmentID *string `json:"rescheduled_from_appointment_id,omitempty"`
	CreatedAt                    *string `json:"created_at,omitempty"`
	UpdatedAt                    *string `json:"updated_at,omitempty"`
}

func fromProtoAppointment(value *recordv1.Appointment) *appointmentResponse {
	if value == nil {
		return nil
	}
	return &appointmentResponse{
		Id:                           value.GetId(),
		StartTime:                    httpdto.Timestamp(value.GetStartTime()),
		EndTime:                      httpdto.Timestamp(value.GetEndTime()),
		Reason:                       value.GetReason(),
		Status:                       httpdto.EnumString(value.GetStatus()),
		PatientID:                    value.GetPatientId(),
		CancelledByUserID:            value.CancelledByUserId,
		RescheduledFromAppointmentID: value.RescheduledFromAppointmentId,
		CreatedAt:                    httpdto.Timestamp(value.GetCreatedAt()),
		UpdatedAt:                    httpdto.Timestamp(value.GetUpdatedAt()),
	}
}

func fromProtoAppointments(values []*recordv1.Appointment) []appointmentResponse {
	if values == nil {
		return nil
	}
	result := make([]appointmentResponse, 0, len(values))
	for _, value := range values {
		mapped := fromProtoAppointment(value)
		if mapped == nil {
			result = append(result, appointmentResponse{})
			continue
		}
		result = append(result, *mapped)
	}
	return result
}
