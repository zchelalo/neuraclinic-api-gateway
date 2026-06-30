package v1

import "github.com/zchelalo/neuraclinic-api-gateway/internal/modules/catalogs/application"

type itemResponse struct {
	Value string `json:"value,omitempty"`
	Label string `json:"label,omitempty"`
}

func fromItems(values []application.Item) []itemResponse {
	if values == nil {
		return nil
	}
	result := make([]itemResponse, 0, len(values))
	for _, value := range values {
		result = append(result, itemResponse{
			Value: value.Value,
			Label: value.Label,
		})
	}
	return result
}
