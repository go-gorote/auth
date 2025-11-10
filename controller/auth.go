package controller

import (
	"github.com/go-gorote/auth/dto"
	"github.com/go-gorote/auth/schema"
	"github.com/go-gorote/auth/secret"
	"github.com/go-gorote/gorote"
	"github.com/gofiber/fiber/v2"
)

// Login godoc
// @Summary      User login
// @Description  Authenticate user with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials body schema.Login true "User login credentials (email and password required)"
// @Success      200 {object} dto.Token "Login successful - returns access_token and refresh_token"
// @Failure      400 {object} dto.ResponseError "Bad request - validation error, invalid body, invalid credentials, or user inactive"
// @Failure      429 {object} dto.ResponseError "Too many requests - rate limit exceeded (60 requests per window)"
// @Router       /auth/login [post]
func (c *AppController) LoginHandler(ctx *fiber.Ctx) error {
	req, ok := ctx.Locals("validatedData").(*schema.Login)
	if !ok {
		c.Logger.ErrorContext(ctx.UserContext(), "invalid login data")
		return fiber.NewError(fiber.StatusBadRequest, "invalid login data")
	}

	c.Logger.InfoContext(ctx.UserContext(), "login attempt",
		"email", req.Email,
		"Host", ctx.Get("Host"),
		"Origin", ctx.Get("Origin"),
		"Content-Type", ctx.Get("Content-Type"),
	)

	user, err := c.Service.Login(req)
	if err != nil {
		c.Logger.ErrorContext(ctx.UserContext(), "login failed", "error", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	accessToken, err := c.Service.GenerateJwt(user, "access_token")
	if err != nil {
		c.Logger.ErrorContext(ctx.UserContext(), "failed to generate access token", "error", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := c.Service.SetCookie(ctx, "access_token", accessToken); err != nil {
		c.Logger.ErrorContext(ctx.UserContext(), "failed to set access token cookie", "error", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	refreshToken, err := c.Service.GenerateJwt(user, "refresh_token")
	if err != nil {
		c.Logger.ErrorContext(ctx.UserContext(), "failed to generate refresh token", "error", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := c.Service.SetCookie(ctx, "refresh_token", refreshToken); err != nil {
		c.Logger.ErrorContext(ctx.UserContext(), "failed to set refresh token cookie", "error", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	c.Logger.InfoContext(ctx.UserContext(), "logged in", "user_id", user.ID.String())

	return ctx.Status(fiber.StatusOK).JSON(dto.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Logout godoc
// @Summary      Logout user
// @Description  Delete access token and refresh token cookies to log out user
// @Tags         Authentication
// @Failure      400 {object} dto.ResponseError "Bad request - validation error, invalid body, invalid credentials, or user inactive"
// @Router       /auth/logout [post]
func (c *AppController) LogoutHandler(ctx *fiber.Ctx) error {
	if err := c.Service.DeleteCookie(ctx, "access_token"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := c.Service.DeleteCookie(ctx, "refresh_token"); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.SendStatus(fiber.StatusOK)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Generate new access token using refresh token (can be sent in body or cookie)
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        refresh_token body schema.RefreshToken true "Refresh token data (optional if sent as cookie)"
// @Success      200 {object} dto.Token "Token refreshed successfully - returns new access_token and same refresh_token"
// @Failure      400 {object} dto.ResponseError "Bad request - validation error, invalid body, user not found, or user inactive"
// @Failure      401 {object} dto.ResponseError "Unauthorized - invalid or expired refresh token"
// @Failure      429 {object} dto.ResponseError "Too many requests - rate limit exceeded (60 requests per window)"
// @Router       /auth/refresh [post]
func (c *AppController) RefreshTokenHandler(ctx *fiber.Ctx) error {
	tokenStr := ctx.Cookies("refresh_token")
	if tokenStr == "" {
		tokenStr = ctx.Locals("validatedData").(*schema.RefreshToken).RefreshToken
	}
	var claims secret.JwtClaims
	if err := c.Service.Claims(&claims, tokenStr); err != nil {
		c.Logger.ErrorContext(ctx.UserContext(), "failed to parse refresh token claims", "error", err)
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	c.Logger.InfoContext(ctx.UserContext(), "refresh attempt", "user_id", claims.ID)

	users, err := c.Service.Users(claims.ID)
	if err != nil {
		c.Logger.ErrorContext(ctx.UserContext(), "failed to refrash token: cannot fetch user", "error", err, "user_id", claims.ID)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if len(users) == 0 {
		c.Logger.ErrorContext(ctx.UserContext(), "failed to refrash token: user not found", "user_id", claims.ID)
		return fiber.NewError(fiber.StatusBadRequest, "id user not found")
	}
	user := users[0]
	if !user.Active {
		c.Logger.ErrorContext(ctx.UserContext(), "failed to refrash token: user inactive", "user_id", user.ID.String())
		return fiber.NewError(fiber.StatusBadRequest, "failed to refrash token: user is inactive")
	}

	if user.UpdatedAt.Unix() > claims.IssuedAt.Unix() {
		c.Logger.ErrorContext(ctx.UserContext(), "failed to refrash token: user updated at is newer than issued at", "user_updated_at", user.UpdatedAt, "issued_at", claims.IssuedAt)
		return fiber.NewError(fiber.StatusBadRequest, "failed to refrash token: user is inactive")
	}

	accessToken, err := c.Service.GenerateJwt(&user, "access_token")
	if err != nil {
		c.Logger.ErrorContext(ctx.UserContext(), "failed to generate access token", "error", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := c.Service.SetCookie(ctx, "access_token", accessToken); err != nil {
		c.Logger.ErrorContext(ctx.UserContext(), "failed to set access token cookie", "error", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	c.Logger.InfoContext(ctx.UserContext(), "refreshed token", "user_id", user.ID.String())

	return ctx.Status(fiber.StatusOK).JSON(dto.Token{
		AccessToken:  accessToken,
		RefreshToken: tokenStr,
	})
}

func (c *AppController) SetCookiePainelAdminHandler(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name:     "VITE_APP_NAME",
		Value:    c.AppName,
		HTTPOnly: false,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "VITE_APP_VERSION",
		Value:    c.AppVersion,
		HTTPOnly: false,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "VITE_APP_API_BASEURL",
		Value:    gorote.MustEnvAsString("VITE_APP_API_BASEURL"),
		HTTPOnly: false,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "VITE_APP_API_DOCUMENTATIONURL",
		Value:    gorote.MustEnvAsString("VITE_APP_API_DOCUMENTATIONURL"),
		HTTPOnly: false,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "VITE_APP_STORAGE_BASEURL",
		Value:    gorote.MustEnvAsString("VITE_APP_STORAGE_BASEURL"),
		HTTPOnly: false,
		Secure:   true,
		SameSite: "None",
		Path:     "/",
	})

	return ctx.Next()
}

