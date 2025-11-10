package controller

import (
	"slices"

	"github.com/go-gorote/auth/dto"
	"github.com/go-gorote/auth/model"
	"github.com/go-gorote/auth/permission"
	"github.com/go-gorote/auth/schema"
	"github.com/go-gorote/auth/secret"
	"github.com/go-gorote/gorote"
	"github.com/gofiber/fiber/v2"
)

// recieveUserHandler godoc
// @Summary      Receive user by id
// @Description  Receives a user by id
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id path string true "Id user"
// @Success      200 {object} dto.UserDto "User received successfully"
// @Failure      400 {object} dto.ResponseError "Failed to receive user"
// @Failure      404 {object} dto.ResponseError "Id user not found"
// @Router       /users/{id} [get]
func (c *AppController) RecieveUserHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*schema.RecieveUser)
	claims := ctx.Locals("claimsData").(*secret.JwtClaims)
	editorPermission := slices.Contains(claims.Permissions, string(permission.PermissionAdmin)) || slices.Contains(claims.Permissions, string(permission.PermissionViewUser)) || claims.IsSuperUser
	editorUser := claims.ID == req.ID

	if !editorPermission {
		if !editorUser {
			return fiber.NewError(fiber.StatusForbidden, "you don't have permission to access this route")
		}
	}

	users, err := c.Service.Users(req.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if len(users) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "id user not found")
	}

	return ctx.Status(fiber.StatusOK).JSON(users[0].ToUserDto())
}

// ListUsersHandler godoc
// @Summary      List users
// @Description  Lists all users registered in the system
// @Tags         User
// @Produce      json
// @Param        page query int false "Page number of users to retrieve"
// @Param        limit query int false "Number of users to retrieve per page"
// @Success      200 {object} dto.ListUsersDto "Users retrieved successfully"
// @Failure      400 {object} dto.ResponseError "Failed to retrieve users"
// @Failure      404 {object} dto.ResponseError "No users found"
// @Router       /users [get]
func (c *AppController) ListUsersHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*schema.Paginate)
	users, err := c.Service.Users()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if len(users) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "no users found")
	}

	for _, user := range users {
		user.Roles = slices.DeleteFunc(user.Roles, func(r model.Role) bool {
			return !r.Active
		})
	}

	countUsers := uint(len(users))
	if err := gorote.Pagination(req.Page, req.Limit, &users); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	var userDtos []dto.UserDto
	for _, user := range users {
		userDtos = append(userDtos, user.ToUserDto())
	}
	res := &dto.ListUsersDto{
		Page:  req.Page,
		Limit: req.Limit,
		Total: countUsers,
		Data:  userDtos,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

// CreateUserHandler godoc
// @Summary      Create a user
// @Description  Creates a new user with profile and avatar
// @Tags         User
// @Accept       multipart/form-data
// @Produce      json
// @Param        email formData string true "User email"
// @Param        username formData string true "Username (min 3 chars, alphanumeric, . or _)"
// @Param        first_name formData string true "First name"
// @Param        last_name formData string false "Last name"
// @Param        active formData boolean false "User active status"
// @Param        is_super_user formData boolean false "User is superuser"
// @Param        roles formData array false "List of roles"
// @Param        tenants formData array false "List of tenants"
// @Param        phone1 formData string true "Primary phone number (E.164 format)"
// @Param        phone2 formData string false "Secondary phone number (E.164 format)"
// @Param        avatar formData file false "Avatar file"
// @Param        password formData string true "Password (8-72 chars)"
// @Success      201 {object} dto.UserDto "User created successfully"
// @Failure      400 {object} dto.ResponseError "Failed to create user"
// @Router       /users [post]
func (c *AppController) CreateUserHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*schema.CreateUser)
	claims := ctx.Locals("claimsData").(*secret.JwtClaims)
	if err := gorote.ValidatePassword(req.Password); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	hashedPassword, err := gorote.HashPassword(req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "crypting password failed")
	}
	user, err := c.Service.CreateUser(ctx, req, hashedPassword, claims.IsSuperUser)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.Status(fiber.StatusCreated).JSON(user.ToUserDto())
}

// UpdateUserHandler godoc
// @Summary      Update a user
// @Description  Updates a user with new data
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id path string true "Id user"
// @Param        req body schema.UpdateUser true "User data"
// @Success      200 {object} dto.UserDto "User updated successfully"
// @Failure      400 {object} dto.ResponseError "Failed to update user"
// @Failure      403 {object} dto.ResponseError "You don't have permission to update this user"
// @Router       /users/{id} [put]
func (c *AppController) UpdateUserHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*schema.UpdateUser)
	claims := ctx.Locals("claimsData").(*secret.JwtClaims)
	editorPermission := slices.Contains(claims.Permissions, string(permission.PermissionAdmin)) || slices.Contains(claims.Permissions, string(permission.PermissionUpdateUser)) || claims.IsSuperUser
	editorUser := claims.ID == req.ID
	var res model.User
	if editorPermission || editorUser || claims.IsSuperUser {
		user, err := c.Service.UpdateUser(req, claims.IsSuperUser, editorPermission)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "failed to update user")
		}
		res = *user
	} else {
		return fiber.NewError(fiber.StatusForbidden, "you don't have permission to update this user")
	}

	return ctx.Status(fiber.StatusOK).JSON(res.ToUserDto())
}
