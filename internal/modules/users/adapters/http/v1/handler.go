package v1

import (
	"net/http"

	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
	userv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/user/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/modules/users/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/server/httpx"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/httpdto"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/meta"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/parse"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
)

const (
	permUserCreate = "PERMISSION_KEY_USER_CREATE"
	permUserView   = "PERMISSION_KEY_USER_VIEW"
	permUserEdit   = "PERMISSION_KEY_USER_EDIT"
	permUserDelete = "PERMISSION_KEY_USER_DELETE"
)

type Handler struct {
	service              *application.Service
	internalServiceToken string
	defaultLimit         int32
}

func NewHandler(service *application.Service, internalServiceToken string, defaultLimit int32) *Handler {
	return &Handler{service: service, internalServiceToken: internalServiceToken, defaultLimit: defaultLimit}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, auth middleware.Middleware, permissions func(...string) middleware.Middleware) {
	mux.Handle("POST /api/v1/users", middleware.Chain(http.HandlerFunc(h.create), auth, permissions(permUserCreate)))
	mux.Handle("GET /api/v1/users/by-email", middleware.Chain(http.HandlerFunc(h.findByEmail), auth, permissions(permUserView)))
	mux.Handle("GET /api/v1/users", middleware.Chain(http.HandlerFunc(h.list), auth, permissions(permUserView)))
	mux.Handle("GET /api/v1/users/{id}", middleware.Chain(http.HandlerFunc(h.findByID), auth, permissions(permUserView)))
	mux.Handle("PATCH /api/v1/users/{id}/password", middleware.Chain(http.HandlerFunc(h.updatePassword), auth))
	mux.Handle("DELETE /api/v1/users/{id}", middleware.Chain(http.HandlerFunc(h.delete), auth, permissions(permUserDelete)))
}

type createUserRequest struct {
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	RoleKey      string    `json:"role_key"`
	Admin        *struct{} `json:"admin,omitempty"`
	Psychologist *struct {
		FirstName      string  `json:"first_name"`
		MiddleName     *string `json:"middle_name,omitempty"`
		FirstLastName  string  `json:"first_last_name"`
		SecondLastName *string `json:"second_last_name,omitempty"`
	} `json:"psychologist,omitempty"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var body createUserRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	roleKey, err := parse.RoleKey(body.RoleKey)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	req := &userv1.UserServiceCreateRequest{
		Email:    body.Email,
		Password: body.Password,
		RoleKey:  roleKey,
	}
	if body.Psychologist != nil {
		req.Profile = &userv1.UserServiceCreateRequest_Psychologist{
			Psychologist: &userv1.PsychologistProfileCreateData{
				FirstName:      body.Psychologist.FirstName,
				MiddleName:     body.Psychologist.MiddleName,
				FirstLastName:  body.Psychologist.FirstLastName,
				SecondLastName: body.Psychologist.SecondLastName,
			},
		}
	} else {
		req.Profile = &userv1.UserServiceCreateRequest_Admin{Admin: &userv1.AdminProfileCreateData{}}
	}
	resp, err := h.service.Create(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), req)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusCreated, fromProtoUserCreateResponse(resp), nil)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	pagination, err := parse.CursorPagination(r, h.defaultLimit)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	req := &userv1.UserServiceListRequest{
		Pagination:  pagination,
		SearchQuery: parse.OptionalString(parse.QueryString(r, "search_query")),
	}
	if roleRaw := parse.QueryString(r, "role_key"); roleRaw != "" {
		roleKey, err := parse.RoleKey(roleRaw)
		if err != nil {
			response.WriteError(w, r, err)
			return
		}
		req.RoleKey = &roleKey
	}
	resp, err := h.service.List(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), req)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoUsers(resp.GetUsers()), meta.FromProto(resp.GetMeta()))
}

func (h *Handler) findByID(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.FindByID(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), r.PathValue("id"))
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoUserFindByIDResponse(resp), nil)
}

func (h *Handler) findByEmail(w http.ResponseWriter, r *http.Request) {
	email := parse.QueryString(r, "email")
	if email == "" {
		response.WriteError(w, r, response.BadRequest("missing_email", "email query param is required", nil))
		return
	}
	resp, err := h.service.FindByEmail(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{
		InternalServiceToken: h.internalServiceToken,
	}), email)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, fromProtoUserFindByEmailResponse(resp), nil)
}

type updatePasswordRequest struct {
	NewPassword string `json:"new_password"`
}

func (h *Handler) updatePassword(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := middleware.RequireCurrentUserOrPermission(r, id, permUserEdit); err != nil {
		response.WriteError(w, r, err)
		return
	}
	var body updatePasswordRequest
	if err := httpx.DecodeBody(r, &body); err != nil {
		response.WriteError(w, r, err)
		return
	}
	resp, err := h.service.UpdatePassword(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{
		InternalServiceToken: h.internalServiceToken,
	}), id, body.NewPassword)
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, httpdto.Operation(resp.GetOperation()), nil)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	resp, err := h.service.Delete(grpcclient.OutgoingContext(r.Context(), r, grpcclient.CallOptions{IncludeAuth: true}), r.PathValue("id"))
	if err != nil {
		response.WriteError(w, r, err)
		return
	}
	response.Write(w, r, http.StatusOK, httpdto.Operation(resp.GetOperation()), nil)
}

var _ = sharedv1.RoleKey_ROLE_KEY_ADMIN
