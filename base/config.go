package base

import (
	"crypto/rsa"
	"time"

	"github.com/go-gorote/gorote/storage"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Config struct {
	*fiber.App
	*gorm.DB
	AppName          string
	AppVersion       string
	PrivateKey       *rsa.PrivateKey
	JwtExpireAccess  time.Duration
	JwtExpireRefresh time.Duration
	SuperEmail       string
	SuperPass        string
	SuperPhone       string
	Domain           string
	Storage          storage.StorageProvider
	Bucket           string
}
