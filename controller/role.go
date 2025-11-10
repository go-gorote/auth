package controller

import (
	"github.com/go-gorote/auth/dto"
	"github.com/go-gorote/auth/schema"
	"github.com/go-gorote/gorote"
	"github.com/gofiber/fiber/v2"
)

// listRolesHandler godoc
// @Summary      List all roles
// @Description  Lists all roles registered in the system
// @Tags         Role
// @Produce      json
// @Param        page query int false "Page number of roles to retrieve"
// @Param        limit query int false "Number of roles to retrieve per page"
// @Success      200 {object} dto.ListRolesDto "Roles retrieved successfully"
// @Failure      400 {object} dto.ResponseError "Failed to retrieve roles"
// @Failure      404 {object} dto.ResponseError "No roles found"
// @Router       /roles [get]
func (c *AppController) ListRolesHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*schema.Paginate)
	roles, err := c.Service.Roles()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if len(roles) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "no roles found")
	}

	countRoles := uint(len(roles))
	if err := gorote.Pagination(req.Page, req.Limit, &roles); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	var data []dto.RoleDto
	for _, role := range roles {
		data = append(data, role.ToRoleDto())
	}
	res := &dto.ListRolesDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: countRoles,
		Data:  data,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

// @Summary      Create a role
// @Description  Creates a role with name, description, active, and permissions
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        req body schema.CreateRole true "Role data"
// @Success      201 {object} dto.RoleDto "Role created successfully"
// @Failure      400 {object} dto.ResponseError "Failed to create role"
// @Router       /roles [post]
func (c *AppController) CreateRoleHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*schema.CreateRole)
	role, err := c.Service.CreateRole(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to create role")
	}
	return ctx.Status(fiber.StatusCreated).JSON(role.ToRoleDto())
}

// UpdateRoleHandler godoc
// @Summary      Update a role
// @Description  Updates a role with new data
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        id path string true "Id user"
// @Param        req body schema.UpdateRole true "Role data"
// @Success      200 {object} dto.RoleDto "Role updated successfully"
// @Failure      400 {object} dto.ResponseError "Failed to update role"
// @Router       /roles/{id} [put]
func (c *AppController) UpdateRoleHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*schema.UpdateRole)
	role, err := c.Service.UpdateRole(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to update role")
	}
	return ctx.Status(fiber.StatusOK).JSON(role.ToRoleDto())
}
