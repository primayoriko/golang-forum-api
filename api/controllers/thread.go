package controllers

import (
	// "fmt"

	"errors"
	"fmt"
	"net/http"
	"strconv"

	validate "github.com/asaskevich/govalidator"
	"gorm.io/gorm"

	"gitlab.com/hydra/forum-api/api/database"
	"gitlab.com/hydra/forum-api/api/models"
	"gitlab.com/hydra/forum-api/api/utils"
)

// GetThreadsList will fetch all threads list of specific criteria
func GetThreadsList(w http.ResponseWriter, r *http.Request) {
	// cUsername := context.Get(r, "username")
	// cUserID := context.Get(r, "id")
	username := r.FormValue("username")

	userIDStr := r.FormValue("userid")
	pageNumStr := r.FormValue("page")
	pageSizeStr := r.FormValue("pagesize")

	fmt.Printf("%s, %s, %s, %s\n", username, userIDStr, pageNumStr, pageSizeStr)
	// if username == "" {
	// 	fmt.Println("yes")
	// }

	if !validate.IsInt(userIDStr) || !validate.IsInt(pageNumStr) || !validate.IsInt(pageSizeStr) {
		utils.JSONResponseWriter(&w, http.StatusBadRequest,
			map[string]interface{}{"message": "bad query value"}, nil)
		return
	}

	userID64, err := strconv.ParseUint(userIDStr, 10, 32)
	userID := uint32(userID64)
	pageNum, err := strconv.Atoi(r.FormValue("page"))
	pageSize, err := strconv.Atoi(r.FormValue("pagesize"))
	offset := pageNum * pageSize

	// if offset < 1 {
	// 	offset := -1
	// }

	db, err := database.ConnectDB()
	if err != nil || db == nil {
		utils.JSONResponseWriter(&w, http.StatusInternalServerError,
			map[string]interface{}{"message": "cannot connect db"}, nil)
		return
	}

	var user models.User
	if username != "" {
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.JSONResponseWriter(&w, http.StatusNotFound,
					nil, nil)
				return
			}

			utils.JSONResponseWriter(&w, http.StatusInternalServerError,
				map[string]interface{}{"message": err}, nil)
			return
		}

		if userID != 0 && userID != user.ID {
			utils.JSONResponseWriter(&w, http.StatusNotFound,
				nil, nil)
			return
		}
	}

	var threads []models.ThreadResult
	if userID != 0 || user.ID != 0 {
		db.Model(&models.Thread{}).Where("ID = ? OR ID = ?", userID, user.ID).
			Offset(offset).
			Limit(pageSize).
			Find(&threads)
	} else {
		db.Model(&models.Thread{}).Offset(offset).
			Limit(pageSize).
			Find(&threads)
	}

	utils.JSONResponseWriter(&w, http.StatusOK,
		map[string]interface{}{"threads": threads}, nil)
	return
}

// GetThread will fetch threads and its posts list of specific criteria
func GetThread(w http.ResponseWriter, r *http.Request) {

}

// CreateThread will make a new thread
func CreateThread(w http.ResponseWriter, r *http.Request) {

}

// UpdateThread will update an existing Thread
func UpdateThread(w http.ResponseWriter, r *http.Request) {

}

// DeleteThread will delete an existing Thread
func DeleteThread(w http.ResponseWriter, r *http.Request) {

}
