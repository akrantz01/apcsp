package websockets

import (
	"sync"
)

// Create mapping of users to client connections
type UserMapping struct {
	sync.RWMutex
	mapping map[string]map[string]*Client
}

// Retrieve all connected clients by user id
func (um *UserMapping) Get(id string) []*Client {
	// Lock for reading
	um.RLock()
	defer um.RUnlock()

	// Ensure mapping exists
	if _, ok := um.mapping[id]; !ok {
		return nil
	}

	// Convert ip-client map to array
	var clients []*Client
	for _, client := range um.mapping[id] {
		clients = append(clients, client)
	}

	// Return mapped clients
	return clients
}

// Map client connection to ip to client id
func (um *UserMapping) Add(id string, client *Client) {
	// Lock for writing
	um.Lock()
	defer um.Unlock()

	// Initialize array of potential clients
	if _, ok := um.mapping[id]; !ok {
		um.mapping[id] = make(map[string]*Client)
	}

	// Add client
	um.mapping[id][client.conn.RemoteAddr().String()] = client
}

func (um *UserMapping) Delete(id string, ip string) {
	// Lock for deletion
	um.Lock()
	defer um.Unlock()

	// Ensure sub-mapping exists
	if _, ok := um.mapping[id]; !ok {
		return
	}

	// Delete the client
	delete(um.mapping[id], ip)
}
