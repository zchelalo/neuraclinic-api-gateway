package v1

import (
	"net/http"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/notes/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/server/httpx"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/httpdto"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/meta"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/parse"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
)

const (
	permNoteCreate = "PERMISSION_KEY_NOTE_CREATE"
	permNoteView   = "PERMISSION_KEY_NOTE_VIEW"
	permNoteEdit   = "PERMISSION_KEY_NOTE_EDIT"
	permNoteDelete = "PERMISSION_KEY_NOTE_DELETE"
)

type Handler struct {
	service      *application.Service
	defaultLimit int32
}

func NewHandler(service *application.Service, defaultLimit int32) *Handler {
	return &Handler{service: service, defaultLimit: defaultLimit}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, auth middleware.Middleware, permissions func(...string) middleware.Middleware) {
	mux.Handle("POST /api/v1/notes", middleware.Chain(http.HandlerFunc(h.create), auth, permissions(permNoteCreate)))
	mux.Handle("GET /api/v1/notes", middleware.Chain(http.HandlerFunc(h.list), auth, permissions(permNoteView)))
	mux.Handle("GET /api/v1/notes/{id}", middleware.Chain(http.HandlerFunc(h.findByID), auth, permissions(permNoteView)))
	mux.Handle("PATCH /api/v1/notes/{id}", middleware.Chain(http.HandlerFunc(h.update), auth, permissions(permNoteEdit)))
	mux.Handle("DELETE /api/v1/notes/{id}", middleware.Chain(http.HandlerFunc(h.delete), auth, permissions(permNoteDelete)))
}

type createRequest struct {
	PatientID     string `json:"patient_id"`
	AppointmentID string `json:"appointment_id"`
	Title         string `json:"title"`
	ContentHTML   string `json:"content_html"`
	ContentText   string `json:"content_text"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var body createRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.Create(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &recordv1.NoteServiceCreateRequest{
		PatientId:     body.PatientID,
		AppointmentId: parse.OptionalString(body.AppointmentID),
		Title:         parse.OptionalString(body.Title),
		ContentHtml:   body.ContentHTML,
		ContentText:   body.ContentText,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusCreated, fromProtoNote(resp.GetNote()), nil)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	pagination, err := parse.CursorPagination(r, h.defaultLimit)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	createdAtRange, err := parse.DateRange(parse.QueryString(r, "created_at_start"), parse.QueryString(r, "created_at_end"))
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	withAppointment, err := parse.QueryOptionalBool(r, "with_appointment_associated")
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	withFiles, err := parse.QueryOptionalBool(r, "with_files_associated")
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.List(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &recordv1.NoteServiceListRequest{
		PatientId:                 parse.QueryString(r, "patient_id"),
		Pagination:                pagination,
		CreatedAtRange:            createdAtRange,
		WithAppointmentAssociated: withAppointment,
		WithFilesAssociated:       withFiles,
		SearchQuery:               parse.OptionalString(parse.QueryString(r, "search_query")),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoNoteSummaries(resp.GetNotes()), meta.FromProto(resp.GetMeta()))
}

func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.FindByID(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), r.PathValue("id"))
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoNote(resp.GetNote()), nil)
}

type updateRequest struct {
	AppointmentID *string `json:"appointment_id,omitempty"`
	Title         *string `json:"title,omitempty"`
	ContentHTML   *string `json:"content_html,omitempty"`
	ContentText   *string `json:"content_text,omitempty"`
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var body updateRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.Update(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &recordv1.NoteServiceUpdateRequest{
		Id:            r.PathValue("id"),
		AppointmentId: body.AppointmentID,
		Title:         body.Title,
		ContentHtml:   body.ContentHTML,
		ContentText:   body.ContentText,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoNote(resp.GetNote()), nil)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.Delete(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), r.PathValue("id"))
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, httpdto.Operation(resp.GetOperation()), nil)
}
