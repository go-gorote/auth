package controller

import (
	"github.com/go-gorote/auth/schema"
	"github.com/gofiber/fiber/v2"
)

// UpdateLogoHandler godoc
// @Summary      Update a logo
// @Description  Updates a logo with new data
// @Tags         Logo
// @Accept       multipart/form-data
// @Produce      json
// @Param        logo formData file true "Logo file"
// @Success      200
// @Failure      400 {object} dto.ResponseError "Failed to update logo"
// @Router       /logo [post]
func (c *AppController) UpdateLogoHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*schema.UpdateLogo)
	if err := c.Service.UpdateLogo(ctx, req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to update logo")
	}
	return ctx.SendStatus(fiber.StatusOK)
}
