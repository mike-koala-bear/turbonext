package main

import "log"

var broadcast = make(chan Message)

func handleMessages() {
	for {
		// Wait for a message to be sent to the broadcast channel
		msg := <-broadcast
		log.Printf("Broadcasting message: %s to %d clients", msg.Content, len(clients)) // Log the number of clients
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error broadcasting message to client: %v", err)
				client.Close()
				delete(clients, client)
			} else {
				log.Printf("Message broadcasted to client")
			}
		}
	}
}
