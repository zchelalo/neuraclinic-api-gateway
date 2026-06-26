package httpx

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/zchelalo/neuraclinic-api-gateway/pkg/response"
)

func DecodeBody[T any](r *http.Request, dst *T) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return response.BadRequest("invalid_json", "invalid json body", err)
	}
	return nil
}

func DecodeBodyAllowEmpty[T any](r *http.Request, dst *T) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return response.BadRequest("invalid_json", "invalid json body", err)
	}
	_ = r.Body.Close()
	if len(bytes.TrimSpace(data)) == 0 {
		return nil
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return response.BadRequest("invalid_json", "invalid json body", err)
	}
	return nil
}
