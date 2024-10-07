package main

import "log"

var broadcast = make(chan Message)

// Handle broadcasting messages to clients
func handleBroadcast() {
	for {
		msg := <-broadcast
		clientsMutex.RLock()
		roomClients, exists := clients[msg.RoomID]
		clientsMutex.RUnlock()

		if !exists {
			continue // No clients in this room
		}

		for client := range roomClients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error writing to WebSocket: %v", err)
				client.Close()

				// Remove the client from the room's client list
				clientsMutex.Lock()
				delete(clients[msg.RoomID], client)
				clientsMutex.Unlock()
			}
		}
	}
}
