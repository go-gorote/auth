package router

import (
	"github.com/go-gorote/auth/permission"
	"github.com/go-gorote/auth/schema"
	"github.com/go-gorote/auth/secret"
	"github.com/go-gorote/gorote"
	"github.com/gofiber/fiber/v2"
)

func (r *AppRouter) UpdateLogo(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.UpdateLogo{}),
			gorote.JWTProtectedRSA(&secret.JwtClaims{}, r.PublicKey, secret.ProtectedRoute(
				permission.PermissionAdmin,
			)),
			r.Controller.UpdateLogoHandler,
		)
	} else {
		h = append(h, handlers...)
	}

	router.Post("/", h...)
}
