package main

import "log"

// Broadcast channel for messages
var Broadcast = make(chan Message)

// Handle broadcasting messages to clients
func handleBroadcast() {
	for {
		msg := <-Broadcast
		clientsMutex.RLock()
		roomClients, exists := clients[msg.RoomID]
		clientsMutex.RUnlock()

		if !exists {
			log.Printf("No clients connected in room %d to broadcast message ID %d", msg.RoomID, msg.ID)
			continue
		}

		log.Printf("Broadcasting message ID %d to room %d with %d clients", msg.ID, msg.RoomID, len(roomClients))

		for client := range roomClients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error writing to WebSocket for room %d: %v", msg.RoomID, err)
				client.Close()

				// Remove the client from the room's client list
				clientsMutex.Lock()
				delete(clients[msg.RoomID], client)
				clientsMutex.Unlock()
				log.Printf("Removed client from room %d due to write error", msg.RoomID)
			}
		}
	}
}
