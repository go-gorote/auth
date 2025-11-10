package seed

import (
	"fmt"

	"github.com/go-gorote/auth/model"
	"github.com/go-gorote/gorote"
	"gorm.io/gorm"
)

func saveUser(db *gorm.DB, user model.User) error {
	if err := gorote.ValidateStruct(user); err != nil {
		return fmt.Errorf("erro de validação")
	}

	if err := gorote.ValidatePassword(user.Password); err != nil {
		return err
	}

	hashPassword, err := gorote.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %s", err.Error())
	}
	user.Password = hashPassword

	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}
