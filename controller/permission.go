package controller

import (
	"github.com/go-gorote/auth/dto"
	"github.com/go-gorote/auth/schema"
	"github.com/go-gorote/gorote"
	"github.com/gofiber/fiber/v2"
)

// @Summary      List all permissions
// @Description  Lists all permissions registered in the system
// @Tags         Permission
// @Produce      json
// @Param        page query int false "Page number of permissions to retrieve"
// @Param        limit query int false "Number of permissions to retrieve per page"
// @Success      200 {object} dto.ListPermissionsDto "Permissions retrieved successfully"
// @Failure      400 {object} dto.ResponseError "Failed to retrieve permissions"
// @Failure      404 {object} dto.ResponseError "No permissions found"
// @Router       /permissions [get]
func (c *AppController) ListPermissiontHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*schema.Paginate)
	permissions, err := c.Service.Permissions()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if len(permissions) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "no permissions found")
	}
	countPermissions := uint(len(permissions))
	if err := gorote.Pagination(req.Page, req.Limit, &permissions); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var data []dto.PermissionDto
	for _, permission := range permissions {
		data = append(data, permission.ToPermissionDto())
	}
	res := &dto.ListPermissionsDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: countPermissions,
		Data:  data,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

// CreatePermissiontHandler godoc
// @Summary      Create a permission
// @Description  Creates a permission with code, description, and active
// @Tags         Permission
// @Accept       json
// @Produce      json
// @Param        req body schema.CreatePermission true "Permission data"
// @Success      201 {object} dto.PermissionDto "Permission created successfully"
// @Failure      400 {object} dto.ResponseError "Failed to create permission"
func (c *AppController) CreatePermissiontHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*schema.CreatePermission)
	permission, err := c.Service.CreatePermission(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to create permission")
	}
	return ctx.Status(fiber.StatusCreated).JSON(permission)
}

// @Summary      Update a permission
// @Description  Updates a permission with new data
// @Tags         Permission
// @Accept       json
// @Produce      json
// @Param        req body schema.UpdatePermission true "Permission data"
// @Success      200 {object} dto.PermissionDto "Permission updated successfully"
// @Failure      400 {object} dto.ResponseError "Failed to update permission"
// @Router       /permissions/{id} [put]
func (c *AppController) UpdatePermissiontHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*schema.UpdatePermission)
	res, err := c.Service.UpdatePermission(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to update permission")
	}
	return ctx.Status(fiber.StatusOK).JSON(res.ToPermissionDto())
}
