package auth

import (
	"fmt"

	"github.com/go-gorote/gorote"
	"gorm.io/gorm"
)

func migrate(config configLoad) error {
	if err := config.db().AutoMigrate(
		&User{},
		&Role{},
		&Permission{},
		&Tenant{},
	); err != nil {
		return err
	}
	return nil
}

func saveUserAdmin(config configLoad) error {
	super := config.super()
	if super.SuperEmail == "" || super.SuperPass == "" {
		return fmt.Errorf("variable super is not set or empty")
	}
	hashPassword, err := gorote.HashPassword(super.SuperPass)
	if err != nil {
		return fmt.Errorf("failed to hash password: %s", err.Error())
	}
	if err := config.db().
		FirstOrCreate(&User{
			Email:       super.SuperEmail,
			Password:    hashPassword,
			Active:      true,
			IsSuperUser: true,
		}).Error; err != nil {
		return err
	}
	return nil
}

func savePermissions(config configLoad) error {
	permissions := []PermissionCode{
		PermissionAdmin,
		PermissionCreateUser,
		PermissionViewUser,
		PermissionUpdateUser,
		PermissionCreatePermission,
		PermissionViewPermission,
		PermissionUpdatePermission,
		PermissionCreateRole,
		PermissionViewRole,
		PermissionUpdateRole,
	}
	for _, permission := range permissions {
		var p Permission
		if err := config.db().Where("code = ?", string(permission)).First(&p).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				p = Permission{Code: string(permission)}
				if err := config.db().Create(&p).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			continue
		}
	}
	return nil
}
