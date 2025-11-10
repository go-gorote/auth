package dto

type PermissionDto struct {
	ID          string `json:"id"`
	UpdatedAt   string `json:"updated_at"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type ListPermissionsDto struct {
	Page  uint            `json:"page"`
	Limit uint            `json:"limit"`
	Total uint            `json:"total"`
	Data  []PermissionDto `json:"data"`
}
