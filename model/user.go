package model

import "github.com/go-gorote/auth/dto"

type User struct {
	BaseModel
	FirstName   string   `gorm:"size:50;not null" validate:"required,min=3,max=50,regexp=^[a-zA-Z]+$" json:"first_name"`
	LastName    string   `gorm:"size:50" validate:"omitempty,max=50,regexp=^[a-zA-Z]+$" json:"last_name"`
	Username    string   `gorm:"uniqueIndex;size:50;not null" validate:"required,min=3,max=50,regexp=^[a-zA-Z0-9._]+$" json:"username"`
	Email       string   `gorm:"uniqueIndex;not null" validate:"required,email" json:"email"`
	Password    string   `gorm:"not null" validate:"required" json:"-"`
	IsSuperUser bool     `gorm:"default:false" json:"is_super_user"`
	Phone1      string   `gorm:"type:varchar(20);not null" validate:"required,e164" json:"phone1"`
	Phone2      string   `gorm:"type:varchar(20)" validate:"omitempty,e164" json:"phone2,omitempty"`
	Roles       []Role   `gorm:"many2many:users_roles" json:"roles"`
	Tenants     []Tenant `gorm:"many2many:users_tenants" json:"tenants"`
	Avatar      string   `json:"avatar"`
	Active      bool     `gorm:"default:true" json:"active"`
}

func (u *User) ToUserDto() dto.UserDto {
	roles := []dto.RoleDto{}
	for _, role := range u.Roles {
		roles = append(roles, role.ToRoleDto())
	}
	tenants := []dto.TenantDto{}
	for _, tenant := range u.Tenants {
		tenants = append(tenants, tenant.ToTenantDto())
	}
	return dto.UserDto{
		ID:          u.ID.String(),
		UpdatedAt:   u.UpdatedAt.Format("02/01/2006 15:04:05"),
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Username:    u.Username,
		Email:       u.Email,
		IsSuperUser: u.IsSuperUser,
		Phone1:      u.Phone1,
		Phone2:      u.Phone2,
		Roles:       roles,
		Tenants:     tenants,
		Avatar:      u.Avatar,
		Active:      u.Active,
	}
}
