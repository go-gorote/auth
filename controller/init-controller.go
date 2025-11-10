package controller

import (
	"log/slog"

	"github.com/go-gorote/auth/service"
	"github.com/gofiber/fiber/v2"
)

type AppController struct {
	AppName    string
	AppVersion string
	Service    service.Service
	Logger     *slog.Logger
}

type Controller interface {
	UpdateLogoHandler(*fiber.Ctx) error
	SetCookiePainelAdminHandler(*fiber.Ctx) error
	// Health
	HealthHandler(*fiber.Ctx) error
	// Auth
	LoginHandler(*fiber.Ctx) error
	LogoutHandler(*fiber.Ctx) error
	RefreshTokenHandler(*fiber.Ctx) error
	// Users
	RecieveUserHandler(*fiber.Ctx) error
	ListUsersHandler(*fiber.Ctx) error
	CreateUserHandler(*fiber.Ctx) error
	UpdateUserHandler(*fiber.Ctx) error
	ChangePasswordHandler(*fiber.Ctx) error
	// Roles
	ListRolesHandler(*fiber.Ctx) error
	CreateRoleHandler(*fiber.Ctx) error
	UpdateRoleHandler(*fiber.Ctx) error
	// Permissions
	ListPermissiontHandler(*fiber.Ctx) error
	CreatePermissiontHandler(*fiber.Ctx) error
	UpdatePermissiontHandler(*fiber.Ctx) error
	// Tenants
	ListTenantHandler(*fiber.Ctx) error
	CreateTenantHandler(*fiber.Ctx) error
	UpdateTenantHandler(*fiber.Ctx) error
}
