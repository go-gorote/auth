package dto

type RoleDto struct {
	ID          string          `json:"id"`
	UpdatedAt   string          `json:"updated_at"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Permissions []PermissionDto `json:"permissions"`
	Active      bool            `json:"active"`
}

type ListRolesDto struct {
	Page  uint      `json:"page"`
	Limit uint      `json:"limit"`
	Total uint      `json:"total"`
	Data  []RoleDto `json:"data"`
}
