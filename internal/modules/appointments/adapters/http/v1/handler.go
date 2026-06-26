package v1

import (
	"net/http"
	"strings"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/appointments/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/server/httpx"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/meta"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/parse"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
)

const (
	permAppointmentCreate = "PERMISSION_KEY_APPOINTMENT_CREATE"
	permAppointmentView   = "PERMISSION_KEY_APPOINTMENT_VIEW"
	permAppointmentEdit   = "PERMISSION_KEY_APPOINTMENT_EDIT"
)

type Handler struct {
	service      *application.Service
	defaultLimit int32
}

func NewHandler(service *application.Service, defaultLimit int32) *Handler {
	return &Handler{service: service, defaultLimit: defaultLimit}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, auth middleware.Middleware, permissions func(...string) middleware.Middleware) {
	mux.Handle("POST /api/v1/appointments", middleware.Chain(http.HandlerFunc(h.create), auth, permissions(permAppointmentCreate)))
	mux.Handle("GET /api/v1/appointments", middleware.Chain(http.HandlerFunc(h.list), auth, permissions(permAppointmentView)))
	mux.Handle("GET /api/v1/appointments/{id}", middleware.Chain(http.HandlerFunc(h.findByID), auth, permissions(permAppointmentView)))
	mux.Handle("PATCH /api/v1/appointments/{id}/reschedule", middleware.Chain(http.HandlerFunc(h.reschedule), auth, permissions(permAppointmentEdit)))
	mux.Handle("PATCH /api/v1/appointments/{id}/status", middleware.Chain(http.HandlerFunc(h.updateStatus), auth, permissions(permAppointmentEdit)))
}

type createRequest struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Reason    string `json:"reason"`
	PatientID string `json:"patient_id"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var body createRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	start, err := parse.Timestamp(body.StartTime)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	end, err := parse.Timestamp(body.EndTime)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.Create(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &recordv1.AppointmentServiceCreateRequest{
		StartTime: start,
		EndTime:   end,
		Reason:    body.Reason,
		PatientId: body.PatientID,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusCreated, fromProtoAppointment(resp.GetAppointment()), nil)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	pagination, err := parse.CursorPagination(r, h.defaultLimit)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	dateRange, err := parse.DateRange(parse.QueryString(r, "start_date"), parse.QueryString(r, "end_date"))
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	statuses, err := parseStatuses(parse.QueryString(r, "statuses"))
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.List(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &recordv1.AppointmentServiceListRequest{
		Pagination: pagination,
		PatientId:  parse.OptionalString(parse.QueryString(r, "patient_id")),
		DateRange:  dateRange,
		Statuses:   statuses,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoAppointments(resp.GetAppointments()), meta.FromProto(resp.GetMeta()))
}

func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.FindByID(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), r.PathValue("id"))
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoAppointment(resp.GetAppointment()), nil)
}

type rescheduleRequest struct {
	NewStartTime string `json:"new_start_time"`
	NewEndTime   string `json:"new_end_time"`
	Reason       string `json:"reason"`
}

func (h *Handler) reschedule(w http.ResponseWriter, r *http.Request) {
	var body rescheduleRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	start, err := parse.Timestamp(body.NewStartTime)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	end, err := parse.Timestamp(body.NewEndTime)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.Reschedule(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &recordv1.AppointmentServiceRescheduleRequest{
		AppointmentId: r.PathValue("id"),
		NewStartTime:  start,
		NewEndTime:    end,
		Reason:        body.Reason,
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoAppointment(resp.GetAppointment()), nil)
}

type updateStatusRequest struct {
	NewStatus         string `json:"new_status"`
	CancelledByUserID string `json:"cancelled_by_user_id"`
}

func (h *Handler) updateStatus(w http.ResponseWriter, r *http.Request) {
	var body updateStatusRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	status, err := parse.AppointmentStatus(body.NewStatus)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.UpdateStatus(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), &recordv1.AppointmentServiceUpdateStatusRequest{
		AppointmentId:     r.PathValue("id"),
		NewStatus:         status,
		CancelledByUserId: parse.OptionalString(body.CancelledByUserID),
	})
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoAppointment(resp.GetAppointment()), nil)
}

func parseStatuses(raw string) ([]sharedv1.AppointmentStatus, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	parts := strings.Split(raw, ",")
	result := make([]sharedv1.AppointmentStatus, 0, len(parts))
	for _, part := range parts {
		value, err := parse.AppointmentStatus(part)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, nil
}
