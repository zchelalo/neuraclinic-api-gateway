package application

import (
	"testing"

	recordv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/record/v1"
	sharedv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/shared/v1"
)

func TestCatalogLocalesCoverAllNonZeroEnumValues(t *testing.T) {
	t.Run("english", func(t *testing.T) {
		assertCatalogCoverage(t, catalogs["en"])
	})
	t.Run("spanish", func(t *testing.T) {
		assertCatalogCoverage(t, catalogs["es"])
	})
}

func TestServiceUsesRequestedLocaleWithEnglishFallback(t *testing.T) {
	service := NewService()

	spanish := service.ListSexes("es-MX,en;q=0.8")
	if len(spanish) == 0 || spanish[0].Label != "Masculino" {
		t.Fatalf("unexpected spanish catalog: %#v", spanish)
	}

	english := service.ListSexes("fr-FR")
	if len(english) == 0 || english[0].Label != "Male" {
		t.Fatalf("unexpected english fallback catalog: %#v", english)
	}
}

func assertCatalogCoverage(t *testing.T, current catalog) {
	t.Helper()

	for _, value := range sexValues {
		assertLabel(t, current.Sexes, value.String())
	}
	for _, value := range maritalStatusValues {
		assertLabel(t, current.MaritalStatuses, value.String())
	}
	for _, value := range appointmentStatusValues {
		assertLabel(t, current.AppointmentStatuses, value.String())
	}
	for _, value := range fileStatusValues {
		assertLabel(t, current.FileStatuses, value.String())
	}
	for _, value := range roleKeyValues {
		assertLabel(t, current.RoleKeys, value.String())
	}

	assertZeroMissing(t, current.Sexes, recordv1.Sex_SEX_UNSPECIFIED.String())
	assertZeroMissing(t, current.MaritalStatuses, recordv1.MaritalStatus_MARITAL_STATUS_UNSPECIFIED.String())
	assertZeroMissing(t, current.AppointmentStatuses, sharedv1.AppointmentStatus_APPOINTMENT_STATUS_UNSPECIFIED.String())
	assertZeroMissing(t, current.FileStatuses, sharedv1.FileStatus_FILE_STATUS_UNSPECIFIED.String())
	assertZeroMissing(t, current.RoleKeys, sharedv1.RoleKey_ROLE_KEY_UNSPECIFIED.String())
}

func assertLabel(t *testing.T, values map[string]string, key string) {
	t.Helper()
	if values[key] == "" {
		t.Fatalf("missing label for %s", key)
	}
}

func assertZeroMissing(t *testing.T, values map[string]string, key string) {
	t.Helper()
	if _, ok := values[key]; ok {
		t.Fatalf("unexpected unspecified enum label for %s", key)
	}
}
