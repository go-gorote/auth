package service

import (
	"fmt"
	"time"

	"github.com/go-gorote/auth/model"
	"github.com/go-gorote/auth/secret"
	"github.com/go-gorote/gorote"
	"github.com/golang-jwt/jwt/v5"
)

func (s *AppService) GenerateJwt(user *model.User, typeToken string) (string, error) {
	var permissions []string
	for _, role := range user.Roles {
		for _, permission := range role.Permissions {
			permissions = append(permissions, permission.Code)
		}
	}
	var tenants []string
	for _, tenant := range user.Tenants {
		tenants = append(tenants, tenant.Name)
	}

	var expire time.Duration
	switch typeToken {
	case "access_token":
		expire = s.JwtExpireAccess
	case "refresh_token":
		expire = s.JwtExpireRefresh
	default:
		return "", fmt.Errorf("invalid token type")
	}

	token, err := gorote.GenerateJwtWithRSA(secret.JwtClaims{
		IsSuperUser: user.IsSuperUser,
		Permissions: permissions,
		Tenants:     tenants,
		Type:        typeToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        user.ID.String(),
			Issuer:    fmt.Sprintf("%s@%s", s.AppName, s.AppVersion),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		},
	}, s.PrivateKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AppService) Claims(claims jwt.Claims, token string) error {
	if err := gorote.ValidateOrGetJWTRSA(claims, token, &s.PrivateKey.PublicKey); err != nil {
		return err
	}
	return nil
}
