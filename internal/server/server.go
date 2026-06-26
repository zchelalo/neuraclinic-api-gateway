package server

import (
	"context"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	appointmentshttp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/appointments/adapters/http/v1"
	appointmentsapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/appointments/application"
	attachmentshttp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/attachments/adapters/http/v1"
	attachmentsapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/attachments/application"
	authhttp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/auth/adapters/http/v1"
	authapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/auth/application"
	familiogramshttp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/familiograms/adapters/http/v1"
	familiogramsapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/familiograms/application"
	locationshttp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/locations/adapters/http/v1"
	locationsapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/locations/application"
	noteshttp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/notes/adapters/http/v1"
	notesapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/notes/application"
	patientshttp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/patients/adapters/http/v1"
	patientsapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/patients/application"
	usershttp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/users/adapters/http/v1"
	usersapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/users/application"
	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
)

type Config struct {
	Port                 int
	ServiceName          string
	AllowedOrigins       []string
	AccessCookieName     string
	RefreshCookieName    string
	CookieDomain         string
	CookieSecure         bool
	InternalServiceToken string
}

type Dependencies struct {
	Auth         *authapp.Service
	Users        *usersapp.Service
	Patients     *patientsapp.Service
	Appointments *appointmentsapp.Service
	Notes        *notesapp.Service
	Attachments  *attachmentsapp.Service
	Familiograms *familiogramsapp.Service
	Locations    *locationsapp.Service
}

type Server struct {
	httpServer  *http.Server
	handler     http.Handler
	serviceName string
	startedAt   time.Time
	metrics     *middleware.AtomicMetrics
}

func New(cfg Config, logger *zap.Logger, deps Dependencies) *Server {
	metrics := &middleware.AtomicMetrics{}
	authMiddleware := middleware.AuthMiddleware(deps.Auth, cfg.AccessCookieName)
	permissionMiddleware := func(required ...string) middleware.Middleware {
		return middleware.PermissionsMiddleware(required...)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /api/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.Write(w, r, http.StatusOK, map[string]any{
			"status":  "ok",
			"service": cfg.ServiceName,
		}, nil)
	}))

	server := &Server{
		serviceName: cfg.ServiceName,
		startedAt:   time.Now().UTC(),
		metrics:     metrics,
	}

	mux.Handle("GET /api/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.Write(w, r, http.StatusOK, map[string]any{
			"service":    server.serviceName,
			"started_at": server.startedAt,
			"uptime_sec": int(time.Since(server.startedAt).Seconds()),
			"requests":   atomic.LoadUint64(&metrics.Requests),
			"errors":     atomic.LoadUint64(&metrics.Errors),
		}, nil)
	}))

	authhttp.NewHandler(deps.Auth, authhttp.CookieConfig{
		Domain:      cfg.CookieDomain,
		Secure:      cfg.CookieSecure,
		AccessName:  cfg.AccessCookieName,
		RefreshName: cfg.RefreshCookieName,
	}).RegisterRoutes(mux, authMiddleware)
	usershttp.NewHandler(deps.Users, cfg.InternalServiceToken, 10).RegisterRoutes(mux, authMiddleware, permissionMiddleware)
	patientshttp.NewHandler(deps.Patients, 10).RegisterRoutes(mux, authMiddleware, permissionMiddleware)
	appointmentshttp.NewHandler(deps.Appointments, 10).RegisterRoutes(mux, authMiddleware, permissionMiddleware)
	noteshttp.NewHandler(deps.Notes, 10).RegisterRoutes(mux, authMiddleware, permissionMiddleware)
	attachmentshttp.NewHandler(deps.Attachments, 10).RegisterRoutes(mux, authMiddleware, permissionMiddleware)
	familiogramshttp.NewHandler(deps.Familiograms).RegisterRoutes(mux, authMiddleware, permissionMiddleware)
	locationshttp.NewHandler(deps.Locations, 20).RegisterRoutes(mux, authMiddleware)

	handler := middleware.Chain(
		mux,
		middleware.AccessControlMiddleware(cfg.AllowedOrigins),
		middleware.RequestIDMiddleware(),
		middleware.LoggerMiddleware(logger, cfg.ServiceName),
		middleware.LogRequestMiddleware(metrics),
	)

	server.httpServer = &http.Server{
		Addr:              listenAddr(cfg.Port),
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}
	server.handler = handler
	return server
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Handler() http.Handler {
	return s.handler
}

func listenAddr(port int) string {
	return ":" + strconv.Itoa(port)
}
