package bootstrap

import (
	"context"
	"fmt"

	appointmentsgrpc "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/appointments/adapters/grpc"
	appointmentsapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/appointments/application"
	attachmentsgrpc "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/attachments/adapters/grpc"
	attachmentsapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/attachments/application"
	authgrpc "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/auth/adapters/grpc"
	authapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/auth/application"
	familiogramsgrpc "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/familiograms/adapters/grpc"
	familiogramsapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/familiograms/application"
	locationsgrpc "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/locations/adapters/grpc"
	locationsapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/locations/application"
	notesgrpc "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/notes/adapters/grpc"
	notesapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/notes/application"
	patientsgrpc "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/patients/adapters/grpc"
	patientsapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/patients/application"
	usersgrpc "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/users/adapters/grpc"
	usersapp "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/users/application"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/server"
	grpcshared "github.com/zchelalo/neuraclinic-api-gateway/internal/shared/grpcclient"
	"go.uber.org/zap"
)

type App struct {
	Server  *server.Server
	Cleanup func(context.Context) error
}

func InitApp(_ context.Context, logger *zap.Logger, cfg Config) (*App, error) {
	clients, err := grpcshared.NewBundle(grpcshared.BundleConfig{
		Auth:           toConnConfig(cfg.AuthGRPC),
		Users:          toConnConfig(cfg.UsersGRPC),
		Records:        toConnConfig(cfg.RecordsGRPC),
		Location:       toConnConfig(cfg.LocationGRPC),
		FileManagement: toConnConfig(cfg.FileManagementGRPC),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot initialize grpc clients: %w", err)
	}

	httpServer := server.New(server.Config{
		Port:                 cfg.Port,
		ServiceName:          cfg.ServiceName,
		AllowedOrigins:       cfg.AllowedOrigins,
		AccessCookieName:     cfg.AccessCookieName,
		RefreshCookieName:    cfg.RefreshCookieName,
		CookieDomain:         cfg.CookieDomain,
		CookieSecure:         cfg.CookieSecure,
		InternalServiceToken: cfg.InternalServiceToken,
	}, logger, server.Dependencies{
		Auth:         authapp.NewService(authgrpc.New(clients.Auth)),
		Users:        usersapp.NewService(usersgrpc.New(clients.Users)),
		Patients:     patientsapp.NewService(patientsgrpc.New(clients.Patients)),
		Appointments: appointmentsapp.NewService(appointmentsgrpc.New(clients.Appointments)),
		Notes:        notesapp.NewService(notesgrpc.New(clients.Notes)),
		Attachments:  attachmentsapp.NewService(attachmentsgrpc.NewRecords(clients.Attachments), attachmentsgrpc.NewFiles(clients.FileManagement)),
		Familiograms: familiogramsapp.NewService(familiogramsgrpc.New(clients.Familiograms)),
		Locations:    locationsapp.NewService(locationsgrpc.New(clients.Locations)),
	})

	return &App{
		Server: httpServer,
		Cleanup: func(ctx context.Context) error {
			_ = httpServer.Shutdown(ctx)
			return clients.Close()
		},
	}, nil
}

func toConnConfig(cfg GRPCConfig) grpcshared.ConnConfig {
	return grpcshared.ConnConfig{
		Addr:               cfg.Addr,
		TLSEnabled:         cfg.TLSEnabled,
		CACertPath:         cfg.CACertPath,
		InsecureSkipVerify: cfg.InsecureSkipVerify,
	}
}
