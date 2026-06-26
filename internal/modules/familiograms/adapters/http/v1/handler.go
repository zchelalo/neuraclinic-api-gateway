package v1

import (
	"net/http"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/familiograms/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/server/httpx"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	permFamiliogramView = "PERMISSION_KEY_PATIENT_VIEW"
	permFamiliogramEdit = "PERMISSION_KEY_PATIENT_EDIT"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, auth middleware.Middleware, permissions func(...string) middleware.Middleware) {
	mux.Handle("GET /api/v1/familiograms/patient/{patient_id}", middleware.Chain(http.HandlerFunc(h.findByPatientID), auth, permissions(permFamiliogramView)))
	mux.Handle("PUT /api/v1/familiograms/patient/{patient_id}", middleware.Chain(http.HandlerFunc(h.updateByPatientID), auth, permissions(permFamiliogramEdit)))
}

func (h *Handler) findByPatientID(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.FindByPatientID(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), r.PathValue("patient_id"))
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoFamiliogram(resp.GetFamiliogram()), nil)
}

type updateRequest struct {
	Data map[string]any `json:"data"`
}

func (h *Handler) updateByPatientID(w http.ResponseWriter, r *http.Request) {
	var body updateRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	data, err := structpb.NewStruct(body.Data)
	if err != nil {
		response.WriteError(w, r, response.BadRequest("invalid_data", "data must be a valid JSON object", err))
		return
	}
	resp, err := h.service.UpdateByPatientID(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), r.PathValue("patient_id"), &recordv1.FamiliogramServiceUpdateRequest{
		Data: data,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoFamiliogram(resp.GetFamiliogram()), nil)
}
