package service

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (s *AppService) SetCookie(ctx *fiber.Ctx, typeToken, token string) error {
	var expire *time.Duration
	switch typeToken {
	case "access_token":
		exp := s.JwtExpireAccess
		expire = &exp
	case "refresh_token":
		exp := s.JwtExpireRefresh
		expire = &exp
	}

	var maxAge int
	if expire != nil {
		maxAge = int(expire.Seconds())
	}

	if isLocalhost(s.Domain) {
		ctx.Cookie(&fiber.Cookie{
			Name:     typeToken,
			Value:    token,
			HTTPOnly: false,
			Secure:   false,
			SameSite: "Lax",
			Path:     "/",
			MaxAge:   maxAge,
		})
	} else {
		ctx.Cookie(&fiber.Cookie{
			Name:     typeToken,
			Value:    token,
			HTTPOnly: true,
			Secure:   true,
			SameSite: "None",
			Path:     "/",
			MaxAge:   maxAge,
			Domain:   s.Domain,
		})
	}

	return nil
}

func (s *AppService) DeleteCookie(ctx *fiber.Ctx, typeToken string) error {
	if isLocalhost(s.Domain) {
		ctx.Cookie(&fiber.Cookie{
			Name:     typeToken,
			Value:    "",
			HTTPOnly: false,
			Secure:   false,
			SameSite: "Lax",
			Path:     "/",
			MaxAge:   -1,
		})
	} else {
		ctx.Cookie(&fiber.Cookie{
			Name:     typeToken,
			Value:    "",
			HTTPOnly: true,
			Secure:   true,
			SameSite: "None",
			Path:     "/",
			MaxAge:   -1,
			Domain:   s.Domain,
		})
	}

	return nil
}

func isLocalhost(domain string) bool {
	return strings.Contains(domain, "localhost") || strings.Contains(domain, "127.0.0.1")
}
