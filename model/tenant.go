package model

import "github.com/go-gorote/auth/dto"

type Tenant struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;size:100;not null" validate:"required,min=3,max=100,regexp=^[a-zA-Z0-9_]+$" json:"name"`
	Description string `json:"description"`
	Url         string `json:"url" validate:"url,omitempty"`
	Logo        string `json:"logo"`
	Active      bool   `json:"active"`
}

func (t Tenant) ToTenantDto() dto.TenantDto {
	return dto.TenantDto{
		ID:          t.ID.String(),
		UpdatedAt:   t.UpdatedAt.Format("02/01/2006 15:04:05"),
		Name:        t.Name,
		Description: t.Description,
		URL:         t.Url,
		Logo:        t.Logo,
		Active:      t.Active,
	}
}
