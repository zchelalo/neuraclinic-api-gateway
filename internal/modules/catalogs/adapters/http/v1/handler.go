package v1

import (
	"net/http"

	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/catalogs/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/language"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
)

type Handler struct {
	service *application.Service
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, auth middleware.Middleware) {
	mux.Handle("GET /api/v1/catalogs/sexes", middleware.Chain(http.HandlerFunc(h.listSexes), auth))
	mux.Handle("GET /api/v1/catalogs/marital-statuses", middleware.Chain(http.HandlerFunc(h.listMaritalStatuses), auth))
	mux.Handle("GET /api/v1/catalogs/appointment-statuses", middleware.Chain(http.HandlerFunc(h.listAppointmentStatuses), auth))
	mux.Handle("GET /api/v1/catalogs/file-statuses", middleware.Chain(http.HandlerFunc(h.listFileStatuses), auth))
	mux.Handle("GET /api/v1/catalogs/role-keys", middleware.Chain(http.HandlerFunc(h.listRoleKeys), auth))
}

func (h *Handler) listSexes(w http.ResponseWriter, r *http.Request) {
	response.Write(w, r, http.StatusOK, fromItems(h.service.ListSexes(language.ResolveRequest(r))), nil)
}

func (h *Handler) listMaritalStatuses(w http.ResponseWriter, r *http.Request) {
	response.Write(w, r, http.StatusOK, fromItems(h.service.ListMaritalStatuses(language.ResolveRequest(r))), nil)
}

func (h *Handler) listAppointmentStatuses(w http.ResponseWriter, r *http.Request) {
	response.Write(w, r, http.StatusOK, fromItems(h.service.ListAppointmentStatuses(language.ResolveRequest(r))), nil)
}

func (h *Handler) listFileStatuses(w http.ResponseWriter, r *http.Request) {
	response.Write(w, r, http.StatusOK, fromItems(h.service.ListFileStatuses(language.ResolveRequest(r))), nil)
}

func (h *Handler) listRoleKeys(w http.ResponseWriter, r *http.Request) {
	response.Write(w, r, http.StatusOK, fromItems(h.service.ListRoleKeys(language.ResolveRequest(r))), nil)
}
