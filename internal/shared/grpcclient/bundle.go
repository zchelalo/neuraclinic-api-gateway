package grpcclient

import (
	"fmt"

	authv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/auth/v1"
	filemanagementv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/file_management/v1"
	locationv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/location/v1"
	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	userv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/user/v1"
)

type BundleConfig struct {
	Auth           ConnConfig
	Users          ConnConfig
	Records        ConnConfig
	Location       ConnConfig
	FileManagement ConnConfig
}

type Bundle struct {
	authConn           closableConn
	usersConn          closableConn
	recordsConn        closableConn
	locationConn       closableConn
	fileManagementConn closableConn

	Auth           authv1.AuthServiceClient
	Users          userv1.UserServiceClient
	Patients       recordv1.PatientServiceClient
	Appointments   recordv1.AppointmentServiceClient
	Notes          recordv1.NoteServiceClient
	Attachments    recordv1.AttachmentServiceClient
	Familiograms   recordv1.FamiliogramServiceClient
	Locations      locationv1.LocationServiceClient
	FileManagement filemanagementv1.FileManagementServiceClient
}

type closableConn interface {
	Close() error
}

func NewBundle(cfg BundleConfig) (*Bundle, error) {
	authConn, err := NewConnection(cfg.Auth)
	if err != nil {
		return nil, fmt.Errorf("init auth connection: %w", err)
	}
	usersConn, err := NewConnection(cfg.Users)
	if err != nil {
		_ = authConn.Close()
		return nil, fmt.Errorf("init users connection: %w", err)
	}
	recordsConn, err := NewConnection(cfg.Records)
	if err != nil {
		_ = authConn.Close()
		_ = usersConn.Close()
		return nil, fmt.Errorf("init records connection: %w", err)
	}
	locationConn, err := NewConnection(cfg.Location)
	if err != nil {
		_ = authConn.Close()
		_ = usersConn.Close()
		_ = recordsConn.Close()
		return nil, fmt.Errorf("init location connection: %w", err)
	}
	fileManagementConn, err := NewConnection(cfg.FileManagement)
	if err != nil {
		_ = authConn.Close()
		_ = usersConn.Close()
		_ = recordsConn.Close()
		_ = locationConn.Close()
		return nil, fmt.Errorf("init file-management connection: %w", err)
	}

	return &Bundle{
		authConn:           authConn,
		usersConn:          usersConn,
		recordsConn:        recordsConn,
		locationConn:       locationConn,
		fileManagementConn: fileManagementConn,
		Auth:               authv1.NewAuthServiceClient(authConn),
		Users:              userv1.NewUserServiceClient(usersConn),
		Patients:           recordv1.NewPatientServiceClient(recordsConn),
		Appointments:       recordv1.NewAppointmentServiceClient(recordsConn),
		Notes:              recordv1.NewNoteServiceClient(recordsConn),
		Attachments:        recordv1.NewAttachmentServiceClient(recordsConn),
		Familiograms:       recordv1.NewFamiliogramServiceClient(recordsConn),
		Locations:          locationv1.NewLocationServiceClient(locationConn),
		FileManagement:     filemanagementv1.NewFileManagementServiceClient(fileManagementConn),
	}, nil
}

func (b *Bundle) Close() error {
	if b == nil {
		return nil
	}
	var firstErr error
	for _, conn := range []closableConn{b.authConn, b.usersConn, b.recordsConn, b.locationConn, b.fileManagementConn} {
		if conn == nil {
			continue
		}
		if err := conn.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
