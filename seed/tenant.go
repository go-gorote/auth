package seed

import (
	"github.com/go-gorote/auth/model"
	"gorm.io/gorm"
)

func SeedTenants(db *gorm.DB) error {
	data := []model.Tenant{
		{
			Name:        "Grupo",
			Description: "Grupo test Grupo",
			Url:         "https://grupo.com.br",
			Active:      true,
		},
		{
			Name:        "Sistema",
			Description: "Sistema test Sistema",
			Url:         "https://sistema.com.br",
			Active:      true,
		},
		{
			Name:        "Gorote",
			Description: "Gorote test Gorote",
			Url:         "https://gorote.com.br",
			Active:      true,
		},
		{
			Name:        "Test",
			Description: "Test test Test",
			Url:         "https://test.com.br",
			Active:      false,
		},
	}

	for i := range data {
		if err := db.Create(&data[i]).Error; err != nil {
			return err
		}
	}
	return nil
}
