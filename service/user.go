package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-gorote/auth/model"
	"github.com/go-gorote/auth/schema"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (s *AppService) Users(ids ...string) ([]model.User, error) {
	var data []model.User
	if len(ids) == 0 {
		if err := s.DB.
			Preload("Roles.Permissions").
			Preload("Tenants").
			Find(&data).Error; err != nil {
			return nil, fmt.Errorf("failed to query database list")
		}
		return data, nil
	}

	if err := s.DB.
		Preload("Roles.Permissions").
		Preload("Tenants").
		Where("id IN ?", ids).
		Find(&data).Error; err != nil {
		return nil, fmt.Errorf("failed to query database")
	}
	return data, nil
}

func (s *AppService) CreateUser(ctx *fiber.Ctx, req *schema.CreateUser, passwordHash string, editorSuper bool) (*model.User, error) {
	var user model.User
	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		user.Email = req.Email
		user.Username = req.Username
		user.Password = passwordHash
		user.FirstName = req.FirstName
		user.LastName = req.LastName
		user.Active = req.Active
		if editorSuper {
			user.IsSuperUser = req.IsSuperUser
		}
		user.Phone1 = req.Phone1
		user.Phone2 = req.Phone2

		if req.Avatar != nil {
			if req.Avatar.Size >= 10*1024 {
				return fmt.Errorf("logo size must be less than 1MB")
			}
			user.Avatar = fmt.Sprintf("uploads-%s/avatar%v%s", s.AppName, time.Now().UnixMicro(), req.Avatar.Filename)

			if s.Storage == nil {
				filePath := fmt.Sprintf("./%s", user.Avatar)
				if err := ctx.SaveFile(req.Avatar, filePath); err != nil {
					return fmt.Errorf("error saving file")
				}
			} else {
				file, err := req.Avatar.Open()
				if err != nil {
					return err
				}
				defer file.Close()
				if err := s.Storage.Upload(c, s.Bucket, user.Avatar, file, req.Avatar.Header.Get("Content-Type")); err != nil {
					return err
				}
			}
		}

		if len(req.Roles) > 0 {
			roles, err := s.Roles(req.Roles...)
			if err != nil {
				return err
			}
			user.Roles = roles
		}
		if len(req.Tenants) > 0 {
			tenants, err := s.Tenants(req.Tenants...)
			if err != nil {
				return err
			}
			user.Tenants = tenants
		}

		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AppService) UpdateUser(req *schema.UpdateUser, editorSuper, editorPermission bool) (*model.User, error) {
	var user model.User
	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		users, err := s.Users(req.ID)
		if err != nil {
			return err
		}
		if len(users) == 0 {
			return fmt.Errorf("no users found")
		}
		user = users[0]

		user.Email = req.Email
		user.Username = req.Username
		user.FirstName = req.FirstName
		user.LastName = req.LastName
		user.Active = req.Active
		user.Phone1 = req.Phone1
		user.Phone2 = req.Phone2
		if editorSuper {
			user.IsSuperUser = req.IsSuperUser
		}
		if err := tx.Model(&user).Select("*").Updates(user).Error; err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		if editorPermission || editorSuper {
			if len(req.Roles) > 0 {
				var roles []model.Role
				if err := tx.
					Preload("Permissions").
					Where("id IN ?", req.Roles).
					Find(&roles).Error; err != nil {
					return fmt.Errorf("failed to fetch roles")
				}
				user.Roles = roles
			} else {
				user.Roles = nil
			}

			if len(req.Tenants) > 0 {
				var tenants []model.Tenant
				if err := tx.
					Where("id IN ?", req.Tenants).
					Find(&tenants).Error; err != nil {
					return fmt.Errorf("failed to fetch tenants")
				}
				user.Tenants = tenants
			} else {
				user.Tenants = nil
			}
		}
		if err := tx.Model(&user).Association("Roles").Replace(user.Roles); err != nil {
			return fmt.Errorf("failed to update roles: %w", err)
		}
		if err := tx.Model(&user).Association("Tenants").Replace(user.Tenants); err != nil {
			return fmt.Errorf("failed to update tenants: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &user, nil
}
