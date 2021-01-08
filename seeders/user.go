package seeders

import (
	"gitlab.com/hydra/forum-api/api/models"
	"gitlab.com/hydra/forum-api/api/utils"
	"gorm.io/gorm"
)

var users []models.User = []models.User{
	models.User{Username: "naufal", Password: utils.HashPasswordNoErr("123"), Email: "p@g.com"},
	models.User{Username: "hasan", Password: utils.HashPasswordNoErr("cp-wf"), Email: "h@g.com"},
	models.User{Username: "taufiq", Password: utils.HashPasswordNoErr("pq-ceo"), Email: "t@g.com"},
	models.User{Username: "dean", Password: utils.HashPasswordNoErr("ctfd"), Email: "d@g.com"},
}

// SeedUsers is function for seed Post data
func SeedUsers(db *gorm.DB) error {
	if tx := db.Create(&users); tx.Error == nil {
		return tx.Error
	}
	return nil
}
