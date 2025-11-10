package controller

import (
	"github.com/go-gorote/auth/dto"
	"github.com/go-gorote/auth/schema"
	"github.com/go-gorote/gorote"
	"github.com/gofiber/fiber/v2"
)

// ListTenantHandler godoc
// @Summary      List all tenants
// @Description  Lists all tenants registered in the system
// @Tags         Tenant
// @Produce      json
// @Param        page query int false "Page number of tenants to retrieve"
// @Param        limit query int false "Number of tenants to retrieve per page"
// @Success      200 {object} dto.ListTenantsDto "Tenants retrieved successfully"
// @Failure      400 {object} dto.ResponseError "Failed to retrieve tenants"
// @Failure      404 {object} dto.ResponseError "No tenants found"
// @Router       /tenants [get]
func (c *AppController) ListTenantHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*schema.Paginate)
	tenants, err := c.Service.Tenants()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if len(tenants) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "no tenants found")
	}
	countPermissions := uint(len(tenants))
	if err := gorote.Pagination(req.Page, req.Limit, &tenants); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	var data []dto.TenantDto
	for _, tenant := range tenants {
		data = append(data, tenant.ToTenantDto())
	}
	res := &dto.ListTenantsDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: countPermissions,
		Data:  data,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

// CreateTenantHandler godoc
// @Summary      Create a tenant
// @Description  Creates a tenant with name and description
// @Tags         Tenant
// @Accept       multipart/form-data
// @Produce      json
// @Param        name formData string true "Tenant name (min 3, max 100 chars)"
// @Param        description formData string false "Tenant description"
// @Param        url formData string false "Tenant website URL"
// @Param        logo formData file false "Tenant logo file"
// @Param        active formData boolean true "Tenant active status"
// @Success      201 {object} dto.TenantDto "Tenant created successfully"
// @Failure      400 {object} dto.ResponseError "Failed to create tenant"
// @Router       /tenants [post]
func (c *AppController) CreateTenantHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*schema.CreateTenant)

	res, err := c.Service.CreateTenant(ctx, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(res.ToTenantDto())
}

// UpdateTenantHandler godoc
// @Summary      Update a tenant
// @Description  Updates a tenant with new data
// @Tags         Tenant
// @Accept       multipart/form-data
// @Produce      json
// @Param        id path string true "Id tenant"
// @Param        name formData string true "Tenant name (min 3, max 100 chars)"
// @Param        description formData string false "Tenant description"
// @Param        url formData string false "Tenant website URL"
// @Param        logo formData file false "Tenant logo file"
// @Param        active formData boolean true "Tenant active status"
// @Success      200 {object} dto.TenantDto "Tenant updated successfully"
// @Failure      400 {object} dto.ResponseError "Failed to update tenant"
// @Router       /tenants/{id} [put]
func (c *AppController) UpdateTenantHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*schema.UpdateTenant)
	res, err := c.Service.UpdateTenant(ctx, req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to update tenant")
	}
	return ctx.Status(fiber.StatusOK).JSON(res.ToTenantDto())
}
