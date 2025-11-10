package controller

import (
	"slices"

	"github.com/go-gorote/auth/permission"
	"github.com/go-gorote/auth/schema"
	"github.com/go-gorote/auth/secret"
	"github.com/gofiber/fiber/v2"
)

// ChangePasswordHandler godoc
// @Summary      Change password
// @Description  Changes the password of a user
// @Tags         Password
// @Accept       json
// @Param        id path string true "Id user"
// @Param        req body schema.ChangePassword true "Password data"
// @Success      200
// @Failure      400 {object} dto.ResponseError "Failed to change password"
// @Router       /users/password/{id} [put]
func (c *AppController) ChangePasswordHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*schema.ChangePassword)
	claims := ctx.Locals("claimsData").(*secret.JwtClaims)
	editorPermission := slices.Contains(claims.Permissions, string(permission.PermissionAdmin)) ||
		slices.Contains(claims.Permissions, string(permission.PermissionUpdateUser)) ||
		claims.IsSuperUser
	editorUser := claims.ID == req.ID
	if editorPermission || editorUser || claims.IsSuperUser {
		if err := c.Service.ChangePassword(req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				"failed to change password: "+err.Error())
		}
	} else {
		return fiber.NewError(fiber.StatusForbidden, "you don't have permission to change password this user")
	}

	return ctx.SendStatus(fiber.StatusOK)
}
