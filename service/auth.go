package service

import (
	"fmt"
	"slices"

	"github.com/go-gorote/auth/model"
	"github.com/go-gorote/auth/schema"
	"github.com/go-gorote/gorote"
)

func (s *AppService) Login(req *schema.Login) (*model.User, error) {
	var user model.User
	result := s.DB.
		Preload("Roles.Permissions").
		Preload("Tenants").
		Where("email = ?", req.Email).
		First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to login: username or password is incorrect")
	}

	if !gorote.CheckPasswordHash(req.Password, user.Password) {
		return nil, fmt.Errorf("failed to login: username or password is incorrect")
	}

	if !user.Active {
		return nil, fmt.Errorf("failed to login: user is inactive")
	}

	user.Tenants = slices.DeleteFunc(user.Tenants, func(t model.Tenant) bool {
		return !t.Active
	})

	user.Roles = slices.DeleteFunc(user.Roles, func(r model.Role) bool {
		return !r.Active
	})

	for _, r := range user.Roles {
		r.Permissions = slices.DeleteFunc(r.Permissions, func(p model.Permission) bool {
			return !p.Active
		})
	}
	return &user, nil
}
