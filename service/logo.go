package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/go-gorote/auth/schema"
	"github.com/gofiber/fiber/v2"
)

func (s *AppService) SetStorageLogo(ctx *fiber.Ctx, f *multipart.FileHeader) error {
	c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if f.Size >= 10*1024*1024 {
		return fmt.Errorf("file size must be less than 10MB")
	}
	data := "uploads/logo.svg"

	if s.Storage == nil {
		filePath := fmt.Sprintf("./%s", data)
		if err := ctx.SaveFile(f, filePath); err != nil {
			return fmt.Errorf("error saving file")
		}
	} else {
		file, err := f.Open()
		if err != nil {
			return err
		}
		defer file.Close()
		if err := s.Storage.Upload(c, s.Bucket, data, file, f.Header.Get("Content-Type")); err != nil {
			return err
		}
	}
	return nil
}

func (s *AppService) UpdateLogo(ctx *fiber.Ctx, req *schema.UpdateLogo) error {
	if req.Logo != nil {
		if err := s.SetStorageLogo(ctx, req.Logo); err != nil {
			return err
		}
	}
	return nil
}
