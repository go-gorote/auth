package dto

type TenantDto struct {
	ID          string `json:"id"`
	UpdatedAt   string `json:"updated_at"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Logo        string `json:"logo"`
	Active      bool   `json:"active"`
}

type ListTenantsDto struct {
	Page  uint        `json:"page"`
	Limit uint        `json:"limit"`
	Total uint        `json:"total"`
	Data  []TenantDto `json:"data"`
}
