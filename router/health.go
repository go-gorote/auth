package router

import "github.com/gofiber/fiber/v2"

func (r *AppRouter) Health(router fiber.Router) {
	router.Get("/", r.Controller.HealthHandler)
}
