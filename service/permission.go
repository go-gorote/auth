package service

import (
	"fmt"

	"github.com/go-gorote/auth/model"
	"github.com/go-gorote/auth/schema"
	"gorm.io/gorm"
)

func (s *AppService) Permissions(ids ...string) ([]model.Permission, error) {
	var permissions []model.Permission
	if len(ids) == 0 {
		if err := s.DB.
			Find(&permissions).Error; err != nil {
			return nil, fmt.Errorf("failed to query database")
		}
		return permissions, nil
	}
	if err := s.DB.
		Where("id IN ?", ids).
		Find(&permissions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch permissions")
	}

	return permissions, nil
}

func (s *AppService) CreatePermission(req *schema.CreatePermission) (*model.Permission, error) {
	var permission model.Permission
	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		permission.Code = req.Code
		permission.Description = req.Description
		permission.Active = req.Active

		if err := tx.Create(&permission).Error; err != nil {
			return fmt.Errorf("failed to create permission")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &permission, nil
}

func (s *AppService) UpdatePermission(req *schema.UpdatePermission) (*model.Permission, error) {
	var permission model.Permission
	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		if req.ID == "" {
			return fmt.Errorf("permission id is required")
		}
		permissions, err := s.Permissions(req.ID)
		if err != nil {
			return err
		}
		if len(permissions) == 0 {
			return fmt.Errorf("no permissions found")
		}
		permission = permissions[0]

		permission.Code = req.Code
		permission.Description = req.Description
		permission.Active = req.Active

		if err := tx.Model(&permission).Select("*").Updates(permission).Error; err != nil {
			return fmt.Errorf("failed to update permission: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &permission, nil
}
