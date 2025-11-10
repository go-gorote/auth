package schema

import (
	"mime/multipart"
)

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

type CreateTenant struct {
	Name        string                `json:"name" validate:"required,min=3,max=100"`
	Description string                `json:"description" validate:"omitempty"`
	Url         string                `json:"url" validate:"url,omitempty"`
	Logo        *multipart.FileHeader `form:"logo" validate:"omitempty"`
	Active      bool                  `json:"active" validate:"omitempty"`
}

type UpdateLogo struct {
	Logo *multipart.FileHeader `form:"logo" validate:"omitempty"`
}

type UpdateTenant struct {
	ID          string                `param:"id" validate:"required"`
	Name        string                `json:"name" validate:"required,min=3,max=100"`
	Description string                `json:"description" validate:"omitempty"`
	Url         string                `json:"url" validate:"url,omitempty"`
	Logo        *multipart.FileHeader `form:"logo" validate:"omitempty"`
	Active      bool                  `json:"active" validate:"omitempty"`
}

type CreateRole struct {
	Name        string   `json:"name" validate:"required,min=3,max=100"`
	Description string   `json:"description" validate:"omitempty"`
	Permissions []string `json:"permissions" validate:"omitempty"`
}

type CreatePermission struct {
	Code        string   `json:"code" validate:"required,min=3,max=100"`
	Description string   `json:"description" validate:"omitempty"`
	Roles       []string `json:"roles" validate:"omitempty"`
	Active      bool     `json:"active" validate:"required"`
}

type UpdatePermission struct {
	ID          string   `param:"id" validate:"required"`
	Code        string   `json:"code" validate:"required,min=3,max=100"`
	Description string   `json:"description" validate:"omitempty"`
	Roles       []string `json:"roles" validate:"omitempty"`
	Active      bool     `json:"active" validate:"omitempty"`
}

type CreateUser struct {
	Email       string                `json:"email" validate:"required,email"`
	Username    string                `json:"username" validate:"required,min=3,max=50,regexp=^[a-zA-Z0-9._]+$"`
	FirstName   string                `json:"first_name" validate:"required,min=1,max=50"`
	LastName    string                `json:"last_name" validate:"omitempty,max=50"`
	Active      bool                  `json:"active" validate:"omitempty"`
	IsSuperUser bool                  `json:"is_super_user" validate:"omitempty"`
	Roles       []string              `json:"roles" validate:"omitempty"`
	Tenants     []string              `json:"tenants" validate:"omitempty"`
	Phone1      string                `json:"phone1" validate:"required,e164"`
	Phone2      string                `json:"phone2" validate:"omitempty,e164"`
	Avatar      *multipart.FileHeader `json:"avatar" validate:"omitempty"`
	Password    string                `json:"password" validate:"required,min=8,max=72"`
}

type RecieveUser struct {
	ID string `param:"id" validate:"required"`
}

type ChangePassword struct {
	ID       string `param:"id" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type UpdateUser struct {
	ID          string   `param:"id" validate:"required"`
	Email       string   `json:"email" validate:"required,email"`
	Username    string   `json:"username" validate:"required,min=3,max=50,regexp=^[a-zA-Z0-9._]+$"`
	FirstName   string   `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string   `json:"last_name" validate:"omitempty,max=50"`
	Active      bool     `json:"active" validate:"omitempty"`
	IsSuperUser bool     `json:"is_super_user" validate:"omitempty"`
	Roles       []string `json:"roles" validate:"omitempty"`
	Tenants     []string `json:"tenants" validate:"omitempty"`
	Phone1      string   `json:"phone1" validate:"required,e164"`
	Phone2      string   `json:"phone2" validate:"omitempty,e164"`
}

type UpdateRole struct {
	ID          string   `param:"id" validate:"required"`
	Name        string   `json:"name" validate:"omitempty,min=1,max=50"`
	Description string   `json:"description" validate:"omitempty,max=50"`
	Permissions []string `json:"permissions" validate:"omitempty"`
	Active      bool     `json:"active" validate:"omitempty"`
}

type Paginate struct {
	Page  uint `query:"page" validate:"required,min=1"`
	Limit uint `query:"limit" validate:"required,min=1"`
}
