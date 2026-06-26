package v1

import (
	authv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/auth/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/middleware"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/httpdto"
)

type tokenResponse struct {
	AccessToken        string  `json:"access_token,omitempty"`
	RefreshToken       *string `json:"refresh_token,omitempty"`
	AccessTokenExpiry  *string `json:"access_token_expiry,omitempty"`
	RefreshTokenExpiry *string `json:"refresh_token_expiry,omitempty"`
}

type webSignInResponse struct {
	SignedIn           bool    `json:"signed_in"`
	AccessTokenExpiry  *string `json:"access_token_expiry"`
	RefreshTokenExpiry *string `json:"refresh_token_expiry"`
}

type webRefreshResponse struct {
	Refreshed           bool    `json:"refreshed"`
	AccessTokenExpiry   *string `json:"access_token_expiry"`
	RefreshTokenExpiry  *string `json:"refresh_token_expiry"`
	RefreshTokenUpdated bool    `json:"refresh_token_updated"`
}

type meResponse struct {
	UserID          string   `json:"user_id"`
	RoleKey         string   `json:"role_key"`
	PsychologistID  string   `json:"psychologist_id"`
	AdminID         string   `json:"admin_id"`
	PermissionsKeys []string `json:"permissions_keys"`
	Mode            string   `json:"mode"`
}

type verifyResetCodeResponse struct {
	ResetToken string `json:"reset_token,omitempty"`
}

type permissionResponse struct {
	Id          string  `json:"id,omitempty"`
	Key         string  `json:"key,omitempty"`
	Description string  `json:"description,omitempty"`
	CreatedAt   *string `json:"created_at,omitempty"`
	UpdatedAt   *string `json:"updated_at,omitempty"`
	DeletedAt   *string `json:"deleted_at,omitempty"`
}

func fromProtoSignInResponse(value *authv1.SignInResponse) *tokenResponse {
	if value == nil {
		return nil
	}
	refreshToken := value.GetRefreshToken()
	return &tokenResponse{
		AccessToken:        value.GetAccessToken(),
		RefreshToken:       optionalString(refreshToken, refreshToken != ""),
		AccessTokenExpiry:  httpdto.Timestamp(value.GetAccessTokenExpiry()),
		RefreshTokenExpiry: httpdto.Timestamp(value.GetRefreshTokenExpiry()),
	}
}

func fromProtoRefreshTokenResponse(value *authv1.RefreshTokenResponse) *tokenResponse {
	if value == nil {
		return nil
	}
	refreshToken := value.GetRefreshToken()
	return &tokenResponse{
		AccessToken:        value.GetAccessToken(),
		RefreshToken:       optionalString(refreshToken, value.RefreshToken != nil),
		AccessTokenExpiry:  httpdto.Timestamp(value.GetAccessTokenExpiry()),
		RefreshTokenExpiry: httpdto.Timestamp(value.GetRefreshTokenExpiry()),
	}
}

func fromAuthContext(value middleware.AuthContext) meResponse {
	return meResponse{
		UserID:          value.UserID,
		RoleKey:         value.RoleKey,
		PsychologistID:  value.PsychologistID,
		AdminID:         value.AdminID,
		PermissionsKeys: value.PermissionsKeys,
		Mode:            string(value.Mode),
	}
}

func fromProtoVerifyResetCodeResponse(value *authv1.VerifyResetCodeResponse) *verifyResetCodeResponse {
	if value == nil {
		return nil
	}
	return &verifyResetCodeResponse{ResetToken: value.GetResetToken()}
}

func fromProtoPermissions(values []*authv1.Permission) []permissionResponse {
	if values == nil {
		return nil
	}
	result := make([]permissionResponse, 0, len(values))
	for _, value := range values {
		result = append(result, permissionResponse{
			Id:          value.GetId(),
			Key:         httpdto.EnumString(value.GetKey()),
			Description: value.GetDescription(),
			CreatedAt:   httpdto.Timestamp(value.GetCreatedAt()),
			UpdatedAt:   httpdto.Timestamp(value.GetUpdatedAt()),
			DeletedAt:   httpdto.Timestamp(value.GetDeletedAt()),
		})
	}
	return result
}

func optionalString(value string, ok bool) *string {
	if !ok {
		return nil
	}
	copied := value
	return &copied
}
