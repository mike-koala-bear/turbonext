"use client" // This ensures it's a client-side component

import { useEffect, useState, useRef } from "react"
import axios from "axios"

export default function ChatPage() {
  const [messages, setMessages] = useState([]) // Store chat messages
  const [input, setInput] = useState("") // Store the user's input
  const socket = useRef(null) // Use ref to avoid recreating WebSocket on each render

  useEffect(() => {
    // Fetch chat history when the component mounts
    axios
      .get("/api/messages")
      .then((response) => {
        setMessages(response.data) // Set chat history (past messages)
      })
      .catch((error) => {
        console.error("Error fetching messages:", error)
      })

    // Function to establish WebSocket connection
    const connectWebSocket = () => {
      socket.current = new WebSocket("ws://localhost:8080/ws") // Update port if needed

      socket.current.onopen = () => {
        console.log("WebSocket connected")
      }

      // Handle incoming WebSocket messages (real-time updates)
      socket.current.onmessage = (event) => {
        const newMessage = JSON.parse(event.data)

        // Prevent duplicate messages: Only add unique messages
        setMessages((prevMessages) => {
          if (!prevMessages.some((msg) => msg.id === newMessage.id)) {
            return [newMessage, ...prevMessages]
          }
          return prevMessages
        })
      }

      // Handle WebSocket closure and attempt reconnection
      socket.current.onclose = (event) => {
        console.log("WebSocket closed:", event)
        if (!event.wasClean) {
          console.log("Attempting to reconnect...")
          setTimeout(connectWebSocket, 3000) // Try to reconnect after 3 seconds
        }
      }

      socket.current.onerror = (error) => {
        console.error("WebSocket error:", error)
      }
    }

    connectWebSocket() // Establish WebSocket connection

    // Cleanup WebSocket when component unmounts
    return () => {
      if (socket.current) {
        console.log("Closing WebSocket")
        socket.current.close()
      }
    }
  }, []) // Empty dependency ensures it only runs once when the component mounts

  // Handle sending a new message
  const handleSendMessage = () => {
    const newMessage = {
      username: "User1", // Replace with dynamic username if needed
      content: input,
    }

    // Send message via WebSocket
    if (socket.current && socket.current.readyState === WebSocket.OPEN) {
      socket.current.send(JSON.stringify(newMessage)) // Send via WebSocket
    }

    // Post message to the backend for persistence (but don't update state manually)
    axios.post("/api/messages", newMessage).catch((error) => {
      console.error("Error posting message:", error)
    })

    setInput("") // Clear the input field after sending the message
  }

  return (
    <div>
      <h1>Chat Room</h1>
      <div className="chat-box">
        {messages.map((msg, index) => (
          <div key={index}>
            <strong>{msg.username}:</strong> {msg.content} <br />
            <small>{new Date(msg.created_at).toLocaleString()}</small>
          </div>
        ))}
      </div>
      <input
        type="text"
        value={input}
        onChange={(e) => setInput(e.target.value)}
        placeholder="Type a message..."
      />
      <button onClick={handleSendMessage}>Send</button>
    </div>
  )
}
