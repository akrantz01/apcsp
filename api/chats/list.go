package chats

import (
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
)

func list(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	logger := logrus.WithFields(logrus.Fields{"app": "chats", "remote_address": r.RemoteAddr, "path": "/api/users", "method": "GET"})

	// Get token w/o validation
	token, err := util.JWT.Unvalidated(r.Header.Get("Authorization"))
	if err != nil {
		logger.WithError(err).Error("Unable to get unvalidated token")
		util.Responses.Error(w, http.StatusInternalServerError, "failed to get token parts")
		return
	}
	logger.Trace("Got unvalidated token parts")

	// Retrieve current user
	var user database.User
	db.Preload("Chats").Preload("Chats.Users").Preload("Chats.Messages").Preload("Chats.Messages.Sender").Preload("Chats.Messages.File").Where("id = ?", util.JWT.Claims(token)["sub"]).First(&user)
	logger.Trace("Got current user's chats and messages")

	// Return empty array if no data
	if len(user.Chats) == 0 {
		util.Responses.SuccessWithData(w, []string{})
		logger.WithFields(logrus.Fields{"chats": len(user.Chats), "id": user.ID}).Debug("Got list of chats for user")
		return
	}

	// Get most recent message for each chat
	for index, chat := range user.Chats {
		user.Chats[index].Messages = []database.Message{chat.Messages[len(chat.Messages)-1]}
	}
	logger.Trace("Set most recent message to only message in chat")

	util.Responses.SuccessWithData(w, user.Chats)
	logger.WithFields(logrus.Fields{"chats": len(user.Chats), "id": user.ID}).Debug("Got list of chats for user")
}
