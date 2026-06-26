package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	authv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/auth/v1"
)

func TestWritePanicsOnProtoPayload(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for proto payload")
		}
	}()

	Write(rec, req, http.StatusOK, &authv1.VerifyResetCodeResponse{ResetToken: "token"}, nil)
}
