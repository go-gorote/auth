package seed

import (
	"github.com/go-gorote/auth/model"
	"gorm.io/gorm"
)

func SeedSuperUser(db *gorm.DB, email, password, phone string) error {
	if err := saveUser(db,
		model.User{
			FirstName:   "Super",
			LastName:    "User",
			Username:    "super.super",
			Email:       email,
			Password:    password,
			IsSuperUser: true,
			Active:      true,
			Phone1:      phone,
		}); err != nil {
		return err
	}
	return nil
}
