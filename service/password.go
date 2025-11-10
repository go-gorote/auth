package service

import (
	"fmt"

	"github.com/go-gorote/auth/schema"
	"github.com/go-gorote/gorote"
)

func (s *AppService) ChangePassword(req *schema.ChangePassword) error {
	users, err := s.Users(req.ID)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return fmt.Errorf("user not found")
	}

	user := users[0]

	if err := gorote.ValidatePassword(req.Password); err != nil {
		return err
	}

	hash, err := gorote.HashPassword(req.Password)
	if err != nil {
		return err
	}

	if err := s.DB.Model(&user).Update("password", hash).Error; err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	return nil
}
