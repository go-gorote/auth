package controller

import "github.com/gofiber/fiber/v2"

// healthHandler godoc
// @Summary      Check service health
// @Description  Returns information about the service health
// @Tags         Health
// @Success      200
// @Failure      500 {object} dto.ResponseError "Service health check failed"
// @Router       /health [get]
func (c *AppController) HealthHandler(ctx *fiber.Ctx) error {
	res, err := c.Service.Health()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}
