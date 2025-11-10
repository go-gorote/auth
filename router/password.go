package router

import (
	"github.com/go-gorote/auth/schema"
	"github.com/go-gorote/auth/secret"
	"github.com/go-gorote/gorote"
	"github.com/gofiber/fiber/v2"
)

func (r *AppRouter) ChangePassword(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.ChangePassword{}),
			gorote.JWTProtectedRSA(&secret.JwtClaims{}, r.PublicKey, secret.ProtectedRoute()),
			r.Controller.ChangePasswordHandler,
		)
	} else {
		h = append(h, handlers...)
	}

	router.Put("/password/:id", h...)
}
