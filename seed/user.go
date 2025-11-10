package seed

import (
	"log"

	"github.com/go-gorote/auth/model"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	data := []model.User{
		{
			FirstName:   "Ralds",
			Username:    "ralds.ralds",
			Email:       "ralds@ralds.com",
			Password:    "Ralds1@#",
			IsSuperUser: true,
			Phone1:      "+5588992200365",
			Active:      true,
		},
		{
			FirstName:   "grupo",
			LastName:    "grupo",
			Username:    "grupo.grupo",
			Email:       "grupo@grupo.com",
			Password:    "Grupo1@#",
			IsSuperUser: false,
			Phone1:      "+5588992200365",
			Active:      false,
		},
		{
			FirstName:   "User",
			LastName:    "User",
			Username:    "user.user",
			Email:       "user@user.com",
			Password:    "Useruser1@#",
			IsSuperUser: false,
			Phone1:      "+5588996877808",
			Active:      true,
		},
		{
			FirstName:   "Sistema",
			LastName:    "Sistema",
			Username:    "sistema.sistema",
			Email:       "sistema@sistema.com",
			Password:    "Sistema1@#",
			IsSuperUser: false,
			Phone1:      "+5588999999999",
			Active:      true,
		},
		{
			FirstName:   "Gorote",
			LastName:    "Gorote",
			Username:    "gorote.gorote",
			Email:       "gorote@gorote.com",
			Password:    "Gorote1@#",
			IsSuperUser: false,
			Phone1:      "+5588999999999",
			Active:      true,
		},
	}

	for i := range data {
		if err := saveUser(db, data[i]); err != nil {
			log.Println(err)
		}
	}
	return nil
}
