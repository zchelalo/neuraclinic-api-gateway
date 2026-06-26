package v1

import (
	"net/http"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/attachments/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/server/httpx"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/httpdto"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/meta"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/parse"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
)

const (
	permFileUpload = "PERMISSION_KEY_FILE_UPLOAD"
	permFileView   = "PERMISSION_KEY_FILE_VIEW"
	permFileDelete = "PERMISSION_KEY_FILE_DELETE"
)

type Handler struct {
	service      *application.Service
	defaultLimit int32
}

func NewHandler(service *application.Service, defaultLimit int32) *Handler {
	return &Handler{service: service, defaultLimit: defaultLimit}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, auth middleware.Middleware, permissions func(...string) middleware.Middleware) {
	mux.Handle("POST /api/v1/attachments", middleware.Chain(http.HandlerFunc(h.create), auth, permissions(permFileUpload)))
	mux.Handle("POST /api/v1/attachments/{id}/confirm-upload", middleware.Chain(http.HandlerFunc(h.confirmUpload), auth, permissions(permFileUpload)))
	mux.Handle("GET /api/v1/attachments", middleware.Chain(http.HandlerFunc(h.list), auth, permissions(permFileView)))
	mux.Handle("GET /api/v1/attachments/{id}", middleware.Chain(http.HandlerFunc(h.findByID), auth, permissions(permFileView)))
	mux.Handle("DELETE /api/v1/attachments/{id}", middleware.Chain(http.HandlerFunc(h.delete), auth, permissions(permFileDelete)))
}

type createRequest struct {
	PatientID    string `json:"patient_id"`
	NoteID       string `json:"note_id"`
	OriginalName string `json:"original_name"`
	MimeType     string `json:"mime_type"`
	SizeBytes    int64  `json:"size_bytes"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var body createRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.Create(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &recordv1.AttachmentServiceCreateRequest{
		PatientId:    body.PatientID,
		NoteId:       parse.OptionalString(body.NoteID),
		OriginalName: body.OriginalName,
		MimeType:     body.MimeType,
		SizeBytes:    body.SizeBytes,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusCreated, createAttachmentResponse{
		Id:        resp.GetId(),
		UploadURL: resp.GetUploadUrl(),
		ExpiresAt: httpdto.Timestamp(resp.GetExpiresAt()),
	}, nil)
}

type confirmUploadRequest struct {
	Status string `json:"status"`
}

func (h *Handler) confirmUpload(w http.ResponseWriter, r *http.Request) {
	var body confirmUploadRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	status, err := parse.FileStatus(body.Status)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	if status != sharedv1.FileStatus_FILE_STATUS_AVAILABLE && status != sharedv1.FileStatus_FILE_STATUS_ERROR {
		response.WriteError(w, r, response.BadRequest("invalid_status", "status must be available or error", nil))
		return
	}
	resp, err := h.service.ConfirmUpload(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), r.PathValue("id"), status)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoConfirmUploadResponse(resp), nil)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	pagination, err := parse.CursorPagination(r, h.defaultLimit)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.List(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &recordv1.AttachmentServiceListRequest{
		Pagination: pagination,
		PatientId:  parse.QueryString(r, "patient_id"),
		NoteId:     parse.OptionalString(parse.QueryString(r, "note_id")),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoAttachments(resp.GetAttachments()), meta.FromProto(resp.GetMeta()))
}

func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.FindByID(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), r.PathValue("id"))
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoAttachment(resp.GetAttachment()), nil)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.Delete(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), r.PathValue("id"))
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, httpdto.Operation(resp.GetOperation()), nil)
}
