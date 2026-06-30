package application

import (
	"embed"
	"encoding/json"
	"fmt"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/language"
)

type Item struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type catalog struct {
	Sexes               map[string]string `json:"sexes"`
	MaritalStatuses     map[string]string `json:"marital_statuses"`
	AppointmentStatuses map[string]string `json:"appointment_statuses"`
	FileStatuses        map[string]string `json:"file_statuses"`
	RoleKeys            map[string]string `json:"role_keys"`
}

//go:embed locales/*.json
var localeFS embed.FS

var catalogs = mustLoadCatalogs()

var (
	sexValues = []recordv1.Sex{
		recordv1.Sex_SEX_MALE,
		recordv1.Sex_SEX_FEMALE,
		recordv1.Sex_SEX_OTHER,
		recordv1.Sex_SEX_PREFER_NOT_TO_SAY,
	}
	maritalStatusValues = []recordv1.MaritalStatus{
		recordv1.MaritalStatus_MARITAL_STATUS_SINGLE,
		recordv1.MaritalStatus_MARITAL_STATUS_MARRIED,
		recordv1.MaritalStatus_MARITAL_STATUS_DIVORCED,
		recordv1.MaritalStatus_MARITAL_STATUS_WIDOWED,
		recordv1.MaritalStatus_MARITAL_STATUS_SEPARATED,
		recordv1.MaritalStatus_MARITAL_STATUS_COHABITING,
	}
	appointmentStatusValues = []sharedv1.AppointmentStatus{
		sharedv1.AppointmentStatus_APPOINTMENT_STATUS_SCHEDULED,
		sharedv1.AppointmentStatus_APPOINTMENT_STATUS_COMPLETED,
		sharedv1.AppointmentStatus_APPOINTMENT_STATUS_RESCHEDULED,
		sharedv1.AppointmentStatus_APPOINTMENT_STATUS_CANCELLED,
		sharedv1.AppointmentStatus_APPOINTMENT_STATUS_NO_SHOW,
	}
	fileStatusValues = []sharedv1.FileStatus{
		sharedv1.FileStatus_FILE_STATUS_UPLOADING,
		sharedv1.FileStatus_FILE_STATUS_AVAILABLE,
		sharedv1.FileStatus_FILE_STATUS_DELETED,
		sharedv1.FileStatus_FILE_STATUS_ERROR,
	}
	roleKeyValues = []sharedv1.RoleKey{
		sharedv1.RoleKey_ROLE_KEY_ADMIN,
		sharedv1.RoleKey_ROLE_KEY_PSYCHOLOGIST,
	}
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ListSexes(locale string) []Item {
	return items(locale, sexValues, func(value recordv1.Sex) string { return value.String() }, func(current catalog) map[string]string {
		return current.Sexes
	})
}

func (s *Service) ListMaritalStatuses(locale string) []Item {
	return items(locale, maritalStatusValues, func(value recordv1.MaritalStatus) string { return value.String() }, func(current catalog) map[string]string {
		return current.MaritalStatuses
	})
}

func (s *Service) ListAppointmentStatuses(locale string) []Item {
	return items(locale, appointmentStatusValues, func(value sharedv1.AppointmentStatus) string { return value.String() }, func(current catalog) map[string]string {
		return current.AppointmentStatuses
	})
}

func (s *Service) ListFileStatuses(locale string) []Item {
	return items(locale, fileStatusValues, func(value sharedv1.FileStatus) string { return value.String() }, func(current catalog) map[string]string {
		return current.FileStatuses
	})
}

func (s *Service) ListRoleKeys(locale string) []Item {
	return items(locale, roleKeyValues, func(value sharedv1.RoleKey) string { return value.String() }, func(current catalog) map[string]string {
		return current.RoleKeys
	})
}

func LocalizedSexLabel(locale string, value recordv1.Sex) string {
	return localizedLabel(locale, value.String(), func(current catalog) map[string]string {
		return current.Sexes
	})
}

func LocalizedMaritalStatusLabel(locale string, value recordv1.MaritalStatus) string {
	return localizedLabel(locale, value.String(), func(current catalog) map[string]string {
		return current.MaritalStatuses
	})
}

func items[T any](locale string, values []T, key func(T) string, labels func(catalog) map[string]string) []Item {
	current := catalogs[normalize(locale)]
	currentLabels := labels(current)
	result := make([]Item, 0, len(values))
	for _, value := range values {
		name := key(value)
		result = append(result, Item{
			Value: name,
			Label: currentLabels[name],
		})
	}
	return result
}

func localizedLabel(locale, fallback string, labels func(catalog) map[string]string) string {
	current := catalogs[normalize(locale)]
	value := labels(current)[fallback]
	if value == "" {
		return fallback
	}
	return value
}

func normalize(locale string) string {
	switch language.ResolveHeader(locale) {
	case language.Spanish:
		return language.Spanish
	default:
		return language.English
	}
}

func mustLoadCatalogs() map[string]catalog {
	return map[string]catalog{
		language.English: mustLoadCatalog(language.English),
		language.Spanish: mustLoadCatalog(language.Spanish),
	}
}

func mustLoadCatalog(locale string) catalog {
	path := fmt.Sprintf("locales/%s.json", locale)
	payload, err := localeFS.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("read catalog locale %s: %v", path, err))
	}
	var current catalog
	if err := json.Unmarshal(payload, &current); err != nil {
		panic(fmt.Sprintf("decode catalog locale %s: %v", path, err))
	}
	if current.Sexes == nil {
		current.Sexes = map[string]string{}
	}
	if current.MaritalStatuses == nil {
		current.MaritalStatuses = map[string]string{}
	}
	if current.AppointmentStatuses == nil {
		current.AppointmentStatuses = map[string]string{}
	}
	if current.FileStatuses == nil {
		current.FileStatuses = map[string]string{}
	}
	if current.RoleKeys == nil {
		current.RoleKeys = map[string]string{}
	}
	return current
}
