package router

import (
	"github.com/go-gorote/auth/permission"
	"github.com/go-gorote/auth/schema"
	"github.com/go-gorote/auth/secret"
	"github.com/go-gorote/gorote"
	"github.com/gofiber/fiber/v2"
)

func (r *AppRouter) ListPermission(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.Paginate{}),
			gorote.JWTProtectedRSA(&secret.JwtClaims{}, r.PublicKey, secret.ProtectedRoute(
				permission.PermissionViewPermission,
				permission.PermissionCreateRole,
				permission.PermissionUpdatePermission,
			)),
			r.Controller.ListPermissiontHandler,
		)
	} else {
		h = append(h, handlers...)
	}

	router.Get("/", h...)
}

func (r *AppRouter) CreatePermission(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.CreatePermission{}),
			gorote.JWTProtectedRSA(&secret.JwtClaims{}, r.PublicKey, secret.ProtectedRoute(
				permission.PermissionCreatePermission,
			)),
			r.Controller.CreatePermissiontHandler,
		)
	} else {
		h = append(h, handlers...)
	}

	router.Post("/", h...)
}

func (r *AppRouter) UpdatePermission(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.UpdatePermission{}),
			gorote.JWTProtectedRSA(&secret.JwtClaims{}, r.PublicKey, secret.ProtectedRoute(
				permission.PermissionUpdatePermission,
			)),
			r.Controller.UpdatePermissiontHandler,
		)
	} else {
		h = append(h, handlers...)
	}

	router.Put("/:id", h...)
}
