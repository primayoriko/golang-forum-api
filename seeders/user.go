package seeders

import (
	"gitlab.com/hydra/forum-api/api/models"
	"gitlab.com/hydra/forum-api/api/utils"
	"gorm.io/gorm"
)

var users []models.User = []models.User{
	{ID: 1, Username: "naufal", Password: utils.HashPasswordNoErr("123"), Email: "p@g.com"},
	{ID: 2, Username: "hasan", Password: utils.HashPasswordNoErr("cp-wf"), Email: "h@g.com"},
	{ID: 3, Username: "taufiq", Password: utils.HashPasswordNoErr("pq-ceo"), Email: "t@g.com"},
	{ID: 4, Username: "dean", Password: utils.HashPasswordNoErr("ctfd"), Email: "d@g.com"}}

// SeedUsers is function for seed Post data
func SeedUsers(db *gorm.DB) error {
	if tx := db.Create(&users); tx.Error == nil {
		return tx.Error
	}
	return nil
}
