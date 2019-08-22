package chats

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"net/http"
)

func list(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Get token w/o validation
	token, err := util.JWT.Unvalidated(r.Header.Get("Authorization"))
	if err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to get token parts")
		return
	}

	// Retrieve current user
	var user database.User
	db.Preload("Chats").Preload("Chats.Users").Preload("Chats.Messages").Preload("Chats.Messages.Sender").Where("id = ?", util.JWT.Claims(token)["sub"]).First(&user)

	// Return empty array if no data
	if len(user.Chats) == 0 {
		util.Responses.SuccessWithData(w, []string{})
		return
	}

	util.Responses.SuccessWithData(w, user.Chats)
}
