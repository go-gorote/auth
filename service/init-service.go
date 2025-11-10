package service

import (
	"log/slog"

	"github.com/go-gorote/auth/base"
	"github.com/go-gorote/auth/model"
	"github.com/go-gorote/auth/schema"
	"github.com/go-gorote/gorote"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AppService struct {
	base.Config
	Logger *slog.Logger
}

type Service interface {
	UpdateLogo(*fiber.Ctx, *schema.UpdateLogo) error
	Health() (*gorote.Health, error)
	SetCookie(*fiber.Ctx, string, string) error
	DeleteCookie(*fiber.Ctx, string) error
	GenerateJwt(*model.User, string) (string, error)
	Login(*schema.Login) (*model.User, error)
	Users(...string) ([]model.User, error)
	Roles(...string) ([]model.Role, error)
	Tenants(...string) ([]model.Tenant, error)
	CreateTenant(*fiber.Ctx, *schema.CreateTenant) (*model.Tenant, error)
	Permissions(...string) ([]model.Permission, error)
	CreatePermission(*schema.CreatePermission) (*model.Permission, error)
	UpdatePermission(*schema.UpdatePermission) (*model.Permission, error)
	CreateRole(*schema.CreateRole) (*model.Role, error)
	CreateUser(*fiber.Ctx, *schema.CreateUser, string, bool) (*model.User, error)
	UpdateUser(*schema.UpdateUser, bool, bool) (*model.User, error)
	UpdateRole(*schema.UpdateRole) (*model.Role, error)
	UpdateTenant(*fiber.Ctx, *schema.UpdateTenant) (*model.Tenant, error)
	ChangePassword(*schema.ChangePassword) error
	Claims(jwt.Claims, string) error
}
