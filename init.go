package auth

import (
	"github.com/go-gorote/auth/base"
	"github.com/go-gorote/auth/controller"
	"github.com/go-gorote/auth/model"
	"github.com/go-gorote/auth/permission"
	"github.com/go-gorote/auth/router"
	"github.com/go-gorote/auth/service"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"gorm.io/gorm"
)

func New(config base.Config) (*router.AppRouter, error) {
	if err := setPermissions(config.DB); err != nil {
		return nil, err
	}

	service := service.AppService{
		Config: config,
		Logger: otelslog.NewLogger("service").With(
			"app_version", config.AppVersion,
			"app_name", config.AppName,
		),
	}

	controller := controller.AppController{
		AppName:    config.AppName,
		AppVersion: config.AppVersion,
		Service:    &service,
		Logger: otelslog.NewLogger("controller").With(
			"app_version", config.AppVersion,
			"app_name", config.AppName,
		),
	}

	router := router.AppRouter{
		App:        config.App,
		PublicKey:  &config.PrivateKey.PublicKey,
		Storage:    config.Storage,
		Controller: &controller,
	}

	return &router, nil
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.Tenant{},
	); err != nil {
		return err
	}
	return nil
}

func setPermissions(db *gorm.DB) error {
	permissions := []permission.PermissionCode{
		// Admin
		permission.PermissionAdmin,
		// Users
		permission.PermissionViewUser,
		permission.PermissionCreateUser,
		permission.PermissionUpdateUser,
		// Permissions
		permission.PermissionViewPermission,
		permission.PermissionCreatePermission,
		permission.PermissionUpdatePermission,
		// Roles
		permission.PermissionViewRole,
		permission.PermissionCreateRole,
		permission.PermissionUpdateRole,
		// Tenants
		permission.PermissionViewTenant,
		permission.PermissionCreateTenant,
		permission.PermissionUpdateTenant,
	}
	for _, permission := range permissions {
		var p model.Permission
		if err := db.Where("code = ?", string(permission)).First(&p).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				p = model.Permission{Code: string(permission), Active: true}
				if err := db.Create(&p).Error; err != nil {
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
