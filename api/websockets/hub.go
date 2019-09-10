package websockets

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
}

// Hub "constructor"
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
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
