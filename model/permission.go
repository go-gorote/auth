package model

import "github.com/go-gorote/auth/dto"

type Permission struct {
	BaseModel
	Code        string `gorm:"uniqueIndex;size:50;not null" validate:"required,regexp=^[a-zA-Z0-9_]+$" json:"code"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

func (p Permission) ToPermissionDto() dto.PermissionDto {
	return dto.PermissionDto{
		ID:          p.ID.String(),
		UpdatedAt:   p.UpdatedAt.Format("02/01/2006 15:04:05"),
		Code:        p.Code,
		Description: p.Description,
		Active:      p.Active,
	}
}
