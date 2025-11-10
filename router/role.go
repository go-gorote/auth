package router

import (
	"github.com/go-gorote/auth/permission"
	"github.com/go-gorote/auth/schema"
	"github.com/go-gorote/auth/secret"
	"github.com/go-gorote/gorote"
	"github.com/gofiber/fiber/v2"
)

func (r *AppRouter) ListRole(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.Paginate{}),
			gorote.JWTProtectedRSA(&secret.JwtClaims{}, r.PublicKey, secret.ProtectedRoute(
				permission.PermissionViewRole,
				permission.PermissionCreateUser,
				permission.PermissionUpdateUser,
			)),
			r.Controller.ListRolesHandler,
		)
	} else {
		h = append(h, handlers...)
	}

	router.Get("/", h...)
}

func (r *AppRouter) CreateRole(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.CreateRole{}),
			gorote.JWTProtectedRSA(&secret.JwtClaims{}, r.PublicKey, secret.ProtectedRoute(
				permission.PermissionCreateRole,
			)),
			r.Controller.CreateRoleHandler,
		)
	} else {
		h = append(h, handlers...)
	}

	router.Post("/", h...)
}

func (r *AppRouter) UpdateRole(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.UpdateRole{}),
			gorote.JWTProtectedRSA(&secret.JwtClaims{}, r.PublicKey, secret.ProtectedRoute(
				permission.PermissionUpdateRole,
			)),
			r.Controller.UpdateRoleHandler,
		)
	} else {
		h = append(h, handlers...)
	}

	router.Put("/:id", h...)
}
