package service

import (
	"fmt"

	"github.com/go-gorote/auth/model"
	"github.com/go-gorote/auth/schema"
	"gorm.io/gorm"
)

func (s *AppService) Roles(ids ...string) ([]model.Role, error) {
	var data []model.Role
	if len(ids) == 0 {
		if err := s.DB.
			Preload("Permissions").
			Find(&data).Error; err != nil {
			return nil, fmt.Errorf("failed to query database")
		}
		return data, nil
	}
	if err := s.DB.
		Preload("Permissions").
		Where("id IN ?", ids).
		Find(&data).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch roles")
	}
	return data, nil
}

func (s *AppService) CreateRole(req *schema.CreateRole) (*model.Role, error) {
	var role model.Role
	if len(req.Permissions) > 0 {
		permissions, err := s.Permissions(req.Permissions...)
		if err != nil {
			return nil, fmt.Errorf("permission with ids does not exist")
		}
		role.Permissions = permissions
	}

	role.Name = req.Name
	role.Description = req.Description
	role.Active = true

	if err := s.DB.Create(&role).Error; err != nil {
		return nil, fmt.Errorf("failed to create role")
	}
	return &role, nil
}

func (s *AppService) UpdateRole(req *schema.UpdateRole) (*model.Role, error) {
	var role model.Role

	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		roles, err := s.Roles(req.ID)
		if err != nil {
			return err
		}
		if len(roles) == 0 {
			return fmt.Errorf("no roles found")
		}
		role = roles[0]

		role.Name = req.Name
		role.Description = req.Description
		role.Active = req.Active

		if len(req.Permissions) > 0 {
			var permissions []model.Permission
			if err := tx.
				Where("id IN ?", req.Permissions).
				Find(&permissions).Error; err != nil {
				return fmt.Errorf("failed to fetch permissions")
			}
			role.Permissions = permissions
		} else {
			role.Permissions = nil
		}

		if err := tx.Model(&role).Select("*").Updates(role).Error; err != nil {
			return fmt.Errorf("failed to update role: %w", err)
		}

		if err := tx.Model(&role).Association("Permissions").Replace(role.Permissions); err != nil {
			return fmt.Errorf("failed to update roles: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &role, nil
}
