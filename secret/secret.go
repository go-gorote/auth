package secret

import (
	"slices"

	"github.com/go-gorote/auth/permission"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	IsSuperUser bool     `json:"isSuperUser"`
	Permissions []string `json:"permissions"`
	Tenants     []string `json:"tenants"`
	Type        string   `json:"type"`
	jwt.RegisteredClaims
}

func ProtectedRoute(p ...permission.PermissionCode) func(jwt.Claims) *fiber.Error {
	return func(c jwt.Claims) *fiber.Error {
		claims := c.(*JwtClaims)
		if claims.Type == "refresh_token" {
			return fiber.NewError(fiber.StatusUnauthorized, "token is refresh token")
		}
		if claims.IsSuperUser {
			return nil
		}
		if len(p) == 0 {
			return nil
		}
		for _, permission := range p {
			if slices.Contains(claims.Permissions, string(permission)) {
				return nil
			}
		}
		return fiber.NewError(fiber.StatusForbidden, "you don't have permission to access this route")
	}
}

func ProtectedRouteWithTenants(tenant *string, p ...permission.PermissionCode) func(jwt.Claims) *fiber.Error {
	return func(c jwt.Claims) *fiber.Error {
		claims := c.(*JwtClaims)
		if claims.Type == "refresh_token" {
			return fiber.NewError(fiber.StatusUnauthorized, "token is refresh token")
		}
		if claims.IsSuperUser {
			return nil
		}

		if tenant != nil {
			if !slices.Contains(claims.Tenants, *tenant) {
				return fiber.NewError(fiber.StatusForbidden, "you don't have permission to access this route")
			}
		}

		if len(p) == 0 {
			return nil
		}
		for _, permission := range p {
			if slices.Contains(claims.Permissions, string(permission)) {
				return nil
			}
		}
		return fiber.NewError(fiber.StatusForbidden, "you don't have permission to access this route")
	}
}
