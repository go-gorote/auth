package router

import (
	"github.com/go-gorote/auth/permission"
	"github.com/go-gorote/auth/schema"
	"github.com/go-gorote/auth/secret"
	"github.com/go-gorote/gorote"
	"github.com/gofiber/fiber/v2"
)

func (r *AppRouter) RecieveUser(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.RecieveUser{}),
			gorote.JWTProtectedRSA(&secret.JwtClaims{}, r.PublicKey, secret.ProtectedRoute()),
			r.Controller.RecieveUserHandler,
		)
	} else {
		h = append(h, handlers...)
	}

	router.Get("/:id", h...)
}

func (r *AppRouter) ListUser(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.Paginate{}),
			gorote.JWTProtectedRSA(&secret.JwtClaims{}, r.PublicKey, secret.ProtectedRoute(
				permission.PermissionViewUser,
				permission.PermissionUpdateUser,
			)),
			r.Controller.ListUsersHandler,
		)
	} else {
		h = append(h, handlers...)
	}

	router.Get("/", h...)
}

func (r *AppRouter) CreateUser(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.CreateUser{}),
			gorote.JWTProtectedRSA(&secret.JwtClaims{}, r.PublicKey, secret.ProtectedRoute(
				permission.PermissionCreateUser,
			)),
			r.Controller.CreateUserHandler,
		)
	} else {
		h = append(h, handlers...)
	}

	router.Post("/", h...)
}

func (r *AppRouter) UpdateUser(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.UpdateUser{}),
			gorote.JWTProtectedRSA(&secret.JwtClaims{}, r.PublicKey, secret.ProtectedRoute()),
			r.Controller.UpdateUserHandler,
		)
	} else {
		h = append(h, handlers...)
	}

	router.Put("/:id", h...)
}
