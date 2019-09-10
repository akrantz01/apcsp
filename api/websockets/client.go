package websockets

import (
	"encoding/json"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/akrantz01/apcsp/api/util"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Websocket configuration
const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// Upgrade connection to websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between websocket connection and the hub
type Client struct {
	hub *Hub

	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// Access to the database
	db *gorm.DB

	// Request logger
	logger *logrus.Entry
}

// readPump pumps messages from the websocket connection to the hub
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine
func (c *Client) readPump() {
	// Close connection and remove from hub
	defer func() {
		c.hub.unregister <- c
		if err := c.conn.Close(); err != nil {
			c.logger.WithError(err).Error("Failed to close websocket connection")
		}
	}()

	// Configure the websocket connection
	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		c.logger.WithError(err).Error("Failed to set read deadline on websocket connection")
	}
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			c.logger.WithError(err).Error("Failed to set read deadline on pong handler")
		}
		return nil
	})

	// Session variables
	authenticated := false
	var user database.User

	for {
		// Read message from connection
		_, rawMsg, err := c.conn.ReadMessage()
		c.logger.Trace("Got new message")

		// Check if connection error is unexpected
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.WithError(err).Error("Unexpected closure of websocket")
			}
			break
		}

		// Parse raw message
		var typeMessage BaseMessage
		if err := json.Unmarshal(rawMsg, &typeMessage); err != nil {
			c.send <- []byte(`{"status": "error", "reason": "unable to decode JSON: ` + err.Error() + `"}`)
			c.logger.WithError(err).Error("Unable to parse json")
			continue
		}
		c.logger.WithField("type", typeMessage.Type).Trace("New message")

		// Ensure authenticated and not authenticating
		if !authenticated && typeMessage.Type != 0 {
			c.send <- []byte(`{"status": "error", "reason": "unauthenticated connection"}`)
			c.logger.Trace("Unauthenticated connection")
			continue
		}

		// Operate on different data based on message type
		switch typeMessage.Type {
		// Handle authentication through websocket
		case MessageAuthentication:
			if authenticated {
				c.logger.Trace("User attempted to re-authenticated")
				continue
			}

			c.logger.Trace("New authentication message")
			var message AuthenticationMessage
			if err := json.Unmarshal(rawMsg, &message); err != nil {
				c.logger.WithError(err).Fatal("Failed to parse json for authentication message. THIS SHOULD NEVER HAPPEN")
				return
			}

			// Validate JWT
			token, err := util.JWT.Validate(message.Token, c.db)
			if err != nil {
				c.logger.WithError(err).Trace("Failed to validate authentication token")
				c.send <- []byte(`{"status": "error", "reason": "invalid token:` + err.Error() + `"}`)
				continue
			}
			c.logger.Trace("Validated authentication token")

			// Get user id from token
			uid, err := util.JWT.UserId(token)
			if err != nil {
				c.logger.WithError(err).Trace("Failed to get user id from token")
				c.send <- []byte(`{"status": "error", "reason": "` + err.Error() + `"}`)
				continue
			}
			c.logger.WithField("uid", uid).Trace("Got user id from token")

			// Set uid in logger
			c.logger = c.logger.WithField("uid", uid)

			// Retrieve user info from database
			c.db.Where("id = ?", uid).First(&user)
			if user.ID == 0 {
				c.logger.Trace("Specified user in token does not exist")
				user = database.User{}
				c.send <- []byte(`{"status": "error", "reason": "user in token does not exist"}`)
				continue
			}
			c.logger.Trace("Retrieved user information from database")

			// Set as authentication
			authenticated = true
			c.logger.Trace("Set connection as authenticated")

			c.send <- []byte(`{"status": "success"}`)
			c.logger.Debug("Authenticated websocket client")

		default:
			c.send <- []byte(`{"status": "error", "reason": "invalid message type"}`)
			c.logger.WithField("type", typeMessage.Type).Info("Invalid message type")
		}
	}
}

// writePump pumps messages from the hub to the websocket connection
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	// Keep connection alive
	// Allows server to know when client dies
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := c.conn.Close(); err != nil {
			c.logger.WithError(err).Error("Failed to close websocket connection")
		}
	}()
	c.logger.Trace("Created new ping-pong ticker")

	for {
		select {
		// Get message to be sent
		case message, ok := <-c.send:
			// Deadlines for messages to be written
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.logger.WithError(err).Error("Failed to set write deadline on websocket connection")
			}
			if !ok {
				// The hub closed the channel
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					c.logger.WithError(err).Error("Failed to set write deadline on websocket control message connection")
				}
				c.logger.Trace("Closed websocket channel")
				return
			}
			c.logger.Trace("Set write deadlines")

			// Get writer for text based messages
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.logger.WithError(err).Trace("Failed to get writer for channel")
				return
			}
			c.logger.Trace("Got new writer for channel")

			// Write the message
			if _, err := w.Write(message); err != nil {
				c.logger.WithError(err).Error("Failed to write message to websocket channel")
			}
			c.logger.Trace("Wrote new message to client")

			// Add queued chat messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				if _, err := w.Write([]byte{'\n'}); err != nil {
					c.logger.WithError(err).Error("Failed to write newline message to websocket channel")
				}
				if _, err := w.Write(<-c.send); err != nil {
					c.logger.WithError(err).Error("Failed to write backlogged message to websocket channel")
				}
			}
			c.logger.Trace("Wrote backlogged messages to client")

			// Close on error
			if err := w.Close(); err != nil {
				c.logger.WithError(err).Trace("Failed to close writer for channel")
				return
			}

		// Handle ping-pongs
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.logger.WithError(err).Error("Failed to set write deadline on ping handler")
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.logger.WithError(err).Trace("Failed to send ping message to client")
				return
			}
		}
	}
}
