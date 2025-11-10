package model

import "github.com/go-gorote/auth/dto"

type Role struct {
	BaseModel
	Name        string       `gorm:"uniqueIndex;size:100;not null" validate:"required,min=3,max=100,regexp=^[a-zA-Z0-9_]+$" json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `gorm:"many2many:roles_permissions" json:"permissions"`
	Active      bool         `json:"active"`
}

func (r *Role) ToRoleDto() dto.RoleDto {
	permissions := []dto.PermissionDto{}

	for _, p := range r.Permissions {
		permissions = append(permissions, p.ToPermissionDto())
	}
	return dto.RoleDto{
		ID:          r.ID.String(),
		UpdatedAt:   r.UpdatedAt.Format("02/01/2006 15:04:05"),
		Name:        r.Name,
		Description: r.Description,
		Permissions: permissions,
		Active:      r.Active,
	}
}
