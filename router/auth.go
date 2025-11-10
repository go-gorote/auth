package router

import (
	"github.com/go-gorote/auth/schema"
	"github.com/go-gorote/gorote"
	"github.com/gofiber/fiber/v2"
)

func (r *AppRouter) Login(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.Login{}),
			r.Controller.LoginHandler,
		)
	} else {
		h = append(h, handlers...)
	}
	router.Post("/login", h...)
}

func (r *AppRouter) Logout(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			r.Controller.LogoutHandler,
		)
	} else {
		h = append(h, handlers...)
	}
	router.Post("/logout", h...)
}

func (r *AppRouter) Refresh(router fiber.Router, handlers ...fiber.Handler) {
	var h []fiber.Handler
	if len(handlers) == 0 {
		h = append(h,
			gorote.ValidationMiddleware(&schema.RefreshToken{}),
			r.Controller.RefreshTokenHandler,
		)
	} else {
		h = append(h, handlers...)
	}

	router.Post("/refresh", h...)
}
