package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/go-gorote/auth/model"
	"github.com/go-gorote/auth/schema"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (s *AppService) Tenants(ids ...string) ([]model.Tenant, error) {
	var data []model.Tenant
	if len(ids) == 0 {
		if err := s.DB.
			Find(&data).Error; err != nil {
			return nil, fmt.Errorf("failed to query database")
		}
		return data, nil
	}
	if err := s.DB.
		Where("id IN ?", ids).
		Find(&data).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch tenants")
	}
	return data, nil
}

func (s *AppService) SetStorage(ctx *fiber.Ctx, f *multipart.FileHeader) (string, error) {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if f.Size >= 10*1024*1024 {
		return "", fmt.Errorf("file size must be less than 10MB")
	}
	ext := filepath.Ext(f.Filename)
	data := fmt.Sprintf("uploads/%v%s", time.Now().UnixMicro(), ext)

	if s.Storage == nil {
		filePath := fmt.Sprintf("./%s", data)
		if err := ctx.SaveFile(f, filePath); err != nil {
			return "", fmt.Errorf("error saving file")
		}
	} else {
		file, err := f.Open()
		if err != nil {
			return "", err
		}
		defer file.Close()
		if err := s.Storage.Upload(c, s.Bucket, data, file, f.Header.Get("Content-Type")); err != nil {
			return "", err
		}
	}
	return data, nil
}

func (s *AppService) CreateTenant(ctx *fiber.Ctx, req *schema.CreateTenant) (*model.Tenant, error) {
	var data model.Tenant

	data.Name = req.Name
	data.Description = req.Description
	data.Url = req.Url
	data.Active = req.Active

	if req.Logo != nil {
		setStorage, err := s.SetStorage(ctx, req.Logo)
		if err != nil {
			return nil, err
		}
		data.Logo = setStorage
	}
	if err := s.DB.Create(&data).Error; err != nil {
		return nil, fmt.Errorf("failed to create role")
	}
	return &data, nil
}

func (s *AppService) UpdateTenant(ctx *fiber.Ctx, req *schema.UpdateTenant) (*model.Tenant, error) {
	var data model.Tenant
	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		if req.ID == "" {
			return fmt.Errorf("tenant id is required")
		}
		tenants, err := s.Tenants(req.ID)
		if err != nil {
			return err
		}
		if len(tenants) == 0 {
			return fmt.Errorf("no tenants found")
		}
		data = tenants[0]

		data.Name = req.Name
		data.Description = req.Description
		data.Url = req.Url
		data.Active = req.Active

		if req.Logo != nil {
			setStorage, err := s.SetStorage(ctx, req.Logo)
			if err != nil {
				return err
			}
			data.Logo = setStorage
		}
		if err := s.DB.Model(&data).Select("*").Updates(data).Error; err != nil {
			return fmt.Errorf("failed to update tenant: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &data, nil
}
