package websockets

import (
	"encoding/json"
	"github.com/akrantz01/apcsp/api/database"
	"github.com/sirupsen/logrus"
	"sync"
)

// Maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Client to user mapping
	mapping UserMapping
}

// Hub "constructor"
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mapping: UserMapping{
			RWMutex: sync.RWMutex{},
			mapping: make(map[string]map[string]*Client),
		},
	}
}

// Main websocket runner
func (h *Hub) Run() {
	// Run forever
	for {
		select {
		// Register a new client
		case client := <-h.register:
			h.clients[client] = true

		// Unregister a client & close connection
		case client := <-h.unregister:
			if _, ok := h.clients[client]; !ok {
				delete(h.clients, client)
				close(client.send)
			}

		// Send message to all clients
		case message := <-h.broadcast:
			// Iterate over all clients
			for client := range h.clients {
				// Send via channel
				select {
				case client.send <- message:

				// If no message, unregister
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// Send message over websocket connection client
func (h *Hub) PushMessage(receiver string, message database.Message, chat string) {
	logger := logrus.WithFields(logrus.Fields{"app": "websocket", "chat": chat, "user": receiver})

	clients := h.mapping.Get(receiver)
	logger.WithField("count", len(clients)).Trace("Got list of clients")

	// Stop if no connections
	if len(clients) == 0 {
		logger.Trace("No clients associated with user")
		return
	}

	// Assemble client message
	msg := ReceiveMessage{
		Type:        MessageReceive,
		Message:     message.Message,
		Chat:        chat,
		Sender:      message.Sender.Username,
		ContentType: int(message.Type),
	}
	logger.Trace("Assembled message to send to clients")

	// Encode message to JSON
	encoded, err := json.Marshal(msg)
	if err != nil {
		logger.WithError(err).Error("Failed to encode to JSON")
		return
	}

	// Send to each client
	for _, client := range clients {
		client.send <- encoded
		logger.Trace("Sent message to client")
	}
}
