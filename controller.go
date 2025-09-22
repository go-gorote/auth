package auth

import (
	"fmt"
	"slices"

	"github.com/go-gorote/gorote"
	"github.com/gofiber/fiber/v2"
)

type controller interface {
	healthHandler(*fiber.Ctx) error
	loginHandler(*fiber.Ctx) error
	refreshTokenHandler(*fiber.Ctx) error
	listUsersHandler(*fiber.Ctx) error
	listPermissiontHandler(*fiber.Ctx) error
	listRolesHandler(*fiber.Ctx) error
	createRoleHandler(*fiber.Ctx) error
	createUserHandler(*fiber.Ctx) error
	updateUserHandler(*fiber.Ctx) error
	recieveUserHandler(*fiber.Ctx) error
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials body login true "User login credentials (email and password required)"
// @Success      200 {object} token "Login successful - returns access_token and refresh_token"
// @Failure      400 {object} map[string]string "Bad request - validation error, invalid body, invalid credentials, or user inactive"
// @Failure      429 {object} map[string]string "Too many requests - rate limit exceeded (60 requests per window)"
// @Router       /auth/login [post]
func (c *appController) loginHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*login)
	user, err := c.service.login(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	accessToken, err := c.service.generateJwt(user, "access_token")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := c.service.setCookie(ctx, "access_token", accessToken); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	refreshToken, err := c.service.generateJwt(user, "refresh_token")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := c.service.setCookie(ctx, "refresh_token", refreshToken); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Generate new access token using refresh token (can be sent in body or cookie)
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        refresh_token body refreshToken true "Refresh token data (optional if sent as cookie)"
// @Success      200 {object} token "Token refreshed successfully - returns new access_token and same refresh_token"
// @Failure      400 {object} map[string]string "Bad request - validation error, invalid body, user not found, or user inactive"
// @Failure      401 {object} map[string]string "Unauthorized - invalid or expired refresh token"
// @Failure      429 {object} map[string]string "Too many requests - rate limit exceeded (60 requests per window)"
// @Router       /auth/refresh [post]
func (c *appController) refreshTokenHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*refreshToken)
	refreshToken := req.RefreshToken
	if refreshToken == "" {
		refreshToken = ctx.Cookies("refresh_token")
	}
	var claims JwtClaims
	if err := c.service.claims(&claims, refreshToken); err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	users, err := c.service.users(claims.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if len(users) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "id user not found")
	}
	user := users[0]
	if !user.Active {
		return fiber.NewError(fiber.StatusBadRequest, "failed to refrash token: user is inactive")
	}

	if user.UpdatedAt.Unix() > claims.IssuedAt.Unix() {
		return fiber.NewError(fiber.StatusBadRequest, "failed to refrash token: user is inactive")
	}

	accessToken, err := c.service.generateJwt(&user, "access_token")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := c.service.setCookie(ctx, "access_token", accessToken); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(token{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
	})
}

// healthHandler godoc
// @Summary      Check service health
// @Description  Returns information about the service health
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{} "Service health information"
// @Failure      500 {object} map[string]string "Service health check failed"
// @Router       /health [get]
func (c *appController) healthHandler(ctx *fiber.Ctx) error {
	res, err := c.service.health()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

// recieveUserHandler godoc
// @Summary      Receive user by id
// @Description  Receives a user by id
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        id path string true "Id user"
// @Success      200 {object} User "User received successfully"
// @Failure      400 {object} map[string]string "Failed to receive user"
// @Failure      404 {object} map[string]string "Id user not found"
// @Router       /receive-user/{id} [get]
func (c *appController) recieveUserHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*recieveUser)

	users, err := c.service.users(req.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if len(users) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "id user not found")
	}
	return ctx.Status(fiber.StatusOK).JSON(users[0])
}

// listUsersHandler godoc
// @Summary      List users
// @Description  Lists all users registered in the system
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number of users to retrieve"
// @Param        limit query int false "Number of users to retrieve per page"
// @Success      200 {object} listUser "Users retrieved successfully"
// @Failure      400 {object} map[string]string "Failed to retrieve users"
// @Failure      404 {object} map[string]string "No users found"
// @Router       /users [get]
func (c *appController) listUsersHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*paginateReq)
	users, err := c.service.users()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if len(users) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "no users found")
	}
	countUsers := uint(len(users))
	if err := gorote.Pagination(req.Page, req.Limit, &users); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	res := &listUser{
		paginateRes: paginateRes{
			Page:  req.Page,
			Limit: req.Limit,
			Total: countUsers,
		},
		Data: users,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

// createUserHandler godoc
// @Summary      Create a user
// @Description  Creates a user with username, password, first name, last name, and phone number
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        req body createUser true "User data"
// @Success      201 {object} map[string]string "User created successfully"
// @Failure      400 {object} map[string]string "Failed to create user"
// @Router       /create-user [post]
func (c *appController) createUserHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*createUser)
	claims := ctx.Locals("claimsData").(*JwtClaims)
	if err := gorote.ValidatePassword(req.Password); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	hashedPassword, err := gorote.HashPassword(req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("crypting password failed: %s", err.Error()))
	}
	req.Password = hashedPassword
	_, err = c.service.createUser(req, claims.IsSuperUser)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.SendStatus(fiber.StatusCreated)
}

// updateUserHandler godoc
// @Summary      Update a user
// @Description  Updates a user with new data
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        req body schemaUser true "User data"
// @Success      200 {object} User "User updated successfully"
// @Failure      400 {object} map[string]string "Failed to update user"
// @Failure      403 {object} map[string]string "You don't have permission to update this user"
// @Router       /update-user [put]
func (c *appController) updateUserHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*schemaUser)
	claims := ctx.Locals("claimsData").(*JwtClaims)
	editorPermission := slices.Contains(claims.Permissions, string(PermissionAdmin)) || slices.Contains(claims.Permissions, string(PermissionUpdateUser))
	editorUser := claims.Subject == req.ID
	var res User
	if editorPermission || editorUser || claims.IsSuperUser {
		user, err := c.service.updateUser(req, claims.IsSuperUser, editorPermission)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "failed to update user")
		}
		res = *user
	} else {
		return fiber.NewError(fiber.StatusForbidden, "you don't have permission to update this user")
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

// listRolesHandler godoc
// @Summary      List all roles
// @Description  Lists all roles registered in the system
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number of roles to retrieve"
// @Param        limit query int false "Number of roles to retrieve per page"
// @Success      200 {object} listRole "Roles retrieved successfully"
// @Failure      400 {object} map[string]string "Failed to retrieve roles"
// @Failure      404 {object} map[string]string "No roles found"
// @Router       /roles [get]
func (c *appController) listRolesHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*paginateReq)
	roles, err := c.service.roles()
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
	res := &listRole{
		paginateRes: paginateRes{
			Page:  req.Page,
			Limit: req.Limit,
			Total: countRoles,
		},
		Data: roles,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

// createRoleHandler godoc
// @Summary      Create a role
// @Description  Creates a role with name, description, and permissions
// @Tags         Role
// @Accept       json
// @Produce      json
// @Param        req body createRole true "Role data"
// @Success      201 {object} Role "Role created successfully"
// @Failure      400 {object} map[string]string "Failed to create role"
// @Router       /create-role [post]
func (c *appController) createRoleHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*createRole)
	role, err := c.service.createRole(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.Status(fiber.StatusCreated).JSON(role)
}

// listPermissiontHandler godoc
// @Summary      List all permissions
// @Description  Lists all permissions registered in the system
// @Tags         Permission
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number of permissions to retrieve"
// @Param        limit query int false "Number of permissions to retrieve per page"
// @Success      200 {object} listPermission "Permissions retrieved successfully"
// @Failure      400 {object} map[string]string "Failed to retrieve permissions"
// @Failure      404 {object} map[string]string "No permissions found"
// @Router       /permissions [get]
func (c *appController) listPermissiontHandler(ctx *fiber.Ctx) error {
	req := ctx.Locals("validatedData").(*paginateReq)
	permissions, err := c.service.permissions()
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
	res := &listPermission{
		paginateRes: paginateRes{
			Page:  req.Page,
			Limit: req.Limit,
			Total: countPermissions,
		},
		Data: permissions,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}
