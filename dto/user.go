package dto

type UserDto struct {
	ID          string      `json:"id"`
	UpdatedAt   string      `json:"updated_at"`
	FirstName   string      `json:"first_name"`
	LastName    string      `json:"last_name"`
	Username    string      `json:"username"`
	Email       string      `json:"email"`
	IsSuperUser bool        `json:"is_super_user"`
	Phone1      string      `json:"phone1"`
	Phone2      string      `json:"phone2,omitempty"`
	Roles       []RoleDto   `json:"roles"`
	Tenants     []TenantDto `json:"tenants"`
	Avatar      string      `json:"avatar"`
	Active      bool        `json:"active"`
}

type ListUsersDto struct {
	Page  uint      `json:"page"`
	Limit uint      `json:"limit"`
	Total uint      `json:"total"`
	Data  []UserDto `json:"data"`
}
