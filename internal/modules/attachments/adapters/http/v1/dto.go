package v1

import (
	filemanagementv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/file_management/v1"
	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/httpdto"
)

type createAttachmentResponse struct {
	Id        string  `json:"id,omitempty"`
	UploadURL string  `json:"upload_url,omitempty"`
	ExpiresAt *string `json:"expires_at,omitempty"`
}

type attachmentResponse struct {
	Id           string  `json:"id,omitempty"`
	FileID       string  `json:"file_id,omitempty"`
	MimeType     string  `json:"mime_type,omitempty"`
	DownloadURL  *string `json:"download_url,omitempty"`
	ExpiresAt    *string `json:"expires_at,omitempty"`
	PatientID    string  `json:"patient_id,omitempty"`
	NoteID       *string `json:"note_id,omitempty"`
	CreatedAt    *string `json:"created_at,omitempty"`
	UpdatedAt    *string `json:"updated_at,omitempty"`
	DeletedAt    *string `json:"deleted_at,omitempty"`
	UploadStatus string  `json:"upload_status,omitempty"`
}

type confirmUploadResponse struct {
	Id          string  `json:"id,omitempty"`
	Status      string  `json:"status,omitempty"`
	DownloadURL *string `json:"download_url,omitempty"`
	ExpiresAt   *string `json:"expires_at,omitempty"`
}

func fromProtoAttachment(value *recordv1.Attachment) *attachmentResponse {
	if value == nil {
		return nil
	}
	return &attachmentResponse{
		Id:           value.GetId(),
		FileID:       value.GetFileId(),
		MimeType:     value.GetMimeType(),
		DownloadURL:  value.DownloadUrl,
		ExpiresAt:    httpdto.Timestamp(value.GetExpiresAt()),
		PatientID:    value.GetPatientId(),
		NoteID:       value.NoteId,
		CreatedAt:    httpdto.Timestamp(value.GetCreatedAt()),
		UpdatedAt:    httpdto.Timestamp(value.GetUpdatedAt()),
		DeletedAt:    httpdto.Timestamp(value.GetDeletedAt()),
		UploadStatus: httpdto.EnumString(value.GetUploadStatus()),
	}
}

func fromProtoAttachments(values []*recordv1.Attachment) []attachmentResponse {
	if values == nil {
		return nil
	}
	result := make([]attachmentResponse, 0, len(values))
	for _, value := range values {
		mapped := fromProtoAttachment(value)
		if mapped == nil {
			result = append(result, attachmentResponse{})
			continue
		}
		result = append(result, *mapped)
	}
	return result
}

func fromProtoConfirmUploadResponse(value *filemanagementv1.FileManagementServiceConfirmUploadResponse) *confirmUploadResponse {
	if value == nil {
		return nil
	}
	return &confirmUploadResponse{
		Id:          value.GetId(),
		Status:      httpdto.EnumString(value.GetStatus()),
		DownloadURL: value.DownloadUrl,
		ExpiresAt:   httpdto.Timestamp(value.GetExpiresAt()),
	}
}
