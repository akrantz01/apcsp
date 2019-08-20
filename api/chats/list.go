package chats

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"net/http"
)

func list(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Validate initial request on headers and body
	if r.Header.Get("Content-Type") != "application/json" {
		util.Responses.Error(w, http.StatusBadRequest, "header 'Content-Type' must be 'application/json'")
		return
	} else if r.Body == nil {
		util.Responses.Error(w, http.StatusBadRequest, "request body must exist")
		return
	}

	// Get token w/o validation
	token, err := util.JWT.Unvalidated(r.Header.Get("Authorization"))
	if err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to get token parts")
		return
	}


	var currentToken database.Token
	db.Where("id = ?", token).First(&currentToken)

	var messages []database.Message
	db.Where("sender_id = ?", currentToken.UserId).Find(&messages)

	if messages == nil {
		util.Responses.Error(w,http.StatusInternalServerError, "Failed to identify any messages from the given user")
		return
	}else{
		util.Responses.SuccessWithData(w, messages)

	}
}

