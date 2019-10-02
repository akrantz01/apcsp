package websockets

import (
	"github.com/akrantz01/apcsp/api/util"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func Websockets(hub *Hub, db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure proper request
		if r.Method != http.MethodGet {
			util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		// Upgrade connection to websockets
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		// Create new client with hub and websocket connection
		client := &Client{
			hub:    hub,
			conn:   conn,
			send:   make(chan []byte, 256),
			db:     db,
			logger: logrus.WithFields(logrus.Fields{"app": "websocket", "remote_address": r.RemoteAddr}),
		}

		// Register with hub
		client.hub.register <- client

		// Start read and write goroutines
		go client.readPump()
		go client.writePump()
	}
}
