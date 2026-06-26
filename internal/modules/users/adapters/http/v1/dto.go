package v1

import (
	userv1 "github.com/zchelalo/neuraclinic-api-gateway/gen/go/user/v1"
	"github.com/zchelalo/neuraclinic-api-gateway/internal/shared/httpdto"
)

type userResponse struct {
	Id        string  `json:"id,omitempty"`
	Email     string  `json:"email,omitempty"`
	RoleKey   string  `json:"role_key,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
	DeletedAt *string `json:"deleted_at,omitempty"`
}

type adminProfileResponse struct {
	Id        string  `json:"id,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
}

type psychologistProfileResponse struct {
	Id             string  `json:"id,omitempty"`
	FirstName      string  `json:"first_name,omitempty"`
	MiddleName     *string `json:"middle_name,omitempty"`
	FirstLastName  string  `json:"first_last_name,omitempty"`
	SecondLastName *string `json:"second_last_name,omitempty"`
	UpdatedAt      *string `json:"updated_at,omitempty"`
}

type userWithProfileResponse struct {
	User         *userResponse                `json:"user,omitempty"`
	Admin        *adminProfileResponse        `json:"admin,omitempty"`
	Psychologist *psychologistProfileResponse `json:"psychologist,omitempty"`
}

func fromProtoUser(value *userv1.User) *userResponse {
	if value == nil {
		return nil
	}
	return &userResponse{
		Id:        value.GetId(),
		Email:     value.GetEmail(),
		RoleKey:   httpdto.EnumString(value.GetRoleKey()),
		CreatedAt: httpdto.Timestamp(value.GetCreatedAt()),
		UpdatedAt: httpdto.Timestamp(value.GetUpdatedAt()),
		DeletedAt: httpdto.Timestamp(value.GetDeletedAt()),
	}
}

func fromProtoUsers(values []*userv1.User) []userResponse {
	if values == nil {
		return nil
	}
	result := make([]userResponse, 0, len(values))
	for _, value := range values {
		mapped := fromProtoUser(value)
		if mapped == nil {
			result = append(result, userResponse{})
			continue
		}
		result = append(result, *mapped)
	}
	return result
}

func fromProtoAdminProfile(value *userv1.AdminProfile) *adminProfileResponse {
	if value == nil {
		return nil
	}
	return &adminProfileResponse{
		Id:        value.GetId(),
		UpdatedAt: httpdto.Timestamp(value.GetUpdatedAt()),
	}
}

func fromProtoPsychologistProfile(value *userv1.PsychologistProfile) *psychologistProfileResponse {
	if value == nil {
		return nil
	}
	return &psychologistProfileResponse{
		Id:             value.GetId(),
		FirstName:      value.GetFirstName(),
		MiddleName:     value.MiddleName,
		FirstLastName:  value.GetFirstLastName(),
		SecondLastName: value.SecondLastName,
		UpdatedAt:      httpdto.Timestamp(value.GetUpdatedAt()),
	}
}

func fromProtoUserCreateResponse(value *userv1.UserServiceCreateResponse) *userWithProfileResponse {
	if value == nil {
		return nil
	}
	return &userWithProfileResponse{
		User:         fromProtoUser(value.GetUser()),
		Admin:        fromProtoAdminProfile(value.GetAdmin()),
		Psychologist: fromProtoPsychologistProfile(value.GetPsychologist()),
	}
}

func fromProtoUserFindByIDResponse(value *userv1.UserServiceFindByIdResponse) *userWithProfileResponse {
	if value == nil {
		return nil
	}
	return &userWithProfileResponse{
		User:         fromProtoUser(value.GetUser()),
		Admin:        fromProtoAdminProfile(value.GetAdmin()),
		Psychologist: fromProtoPsychologistProfile(value.GetPsychologist()),
	}
}

func fromProtoUserFindByEmailResponse(value *userv1.UserServiceFindByEmailResponse) *userWithProfileResponse {
	if value == nil {
		return nil
	}
	return &userWithProfileResponse{
		User:         fromProtoUser(value.GetUser()),
		Admin:        fromProtoAdminProfile(value.GetAdmin()),
		Psychologist: fromProtoPsychologistProfile(value.GetPsychologist()),
	}
}
