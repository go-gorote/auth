package router

import (
	"github.com/go-gorote/auth/permission"
	"github.com/go-gorote/auth/schema"
	"github.com/go-gorote/auth/secret"
	"github.com/go-gorote/gorote"
	"github.com/gofiber/fiber/v2"
)

func (r *AppRouter) ListTenant(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.Paginate{}),
			gorote.JWTProtectedRSA(&secret.JwtClaims{}, r.PublicKey, secret.ProtectedRoute(
				permission.PermissionViewTenant,
				permission.PermissionCreateUser,
				permission.PermissionUpdateUser,
			)),
			r.Controller.ListTenantHandler,
		)
	} else {
		h = append(h, handlers...)
	}
	router.Get("/", h...)
}

func (r *AppRouter) CreateTenant(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.CreateTenant{}),
			gorote.JWTProtectedRSA(&secret.JwtClaims{}, r.PublicKey, secret.ProtectedRoute(
				permission.PermissionCreateTenant,
			)),
			r.Controller.CreateTenantHandler,
		)
	} else {
		h = append(h, handlers...)
	}

	router.Post("/", h...)
}

func (r *AppRouter) UpdateTenant(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.UpdateTenant{}),
			gorote.JWTProtectedRSA(&secret.JwtClaims{}, r.PublicKey, secret.ProtectedRoute(
				permission.PermissionUpdateTenant,
			)),
			r.Controller.UpdateTenantHandler,
		)
	} else {
		h = append(h, handlers...)
	}

	router.Put("/:id", h...)
}
