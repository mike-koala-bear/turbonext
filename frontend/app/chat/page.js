"use client" // Ensure this is a client-side component

import { useEffect, useState, useRef } from "react"
import axios from "axios"

export default function ChatPage() {
  const [messages, setMessages] = useState([]) // Store chat messages
  const [input, setInput] = useState("") // Store the user's input
  const socket = useRef(null) // Use ref to avoid recreating WebSocket on each render

  // Fetch chat history and set up WebSocket connection
  useEffect(() => {
    // Fetch initial chat messages from the API
    axios
      .get("/api/messages")
      .then((response) => {
        setMessages(response.data) // Set messages on initial load
      })
      .catch((error) => {
        console.error("Error fetching messages:", error)
      })

    // Open WebSocket connection to handle real-time updates
    socket.current = new WebSocket("ws://localhost:8080/ws")

    // Handle WebSocket messages (real-time message broadcast)
    socket.current.onmessage = (event) => {
      const newMessage = JSON.parse(event.data)
      // Append new message via WebSocket without re-fetching from API
      setMessages((prevMessages) => [newMessage, ...prevMessages])
    }

    // Clean up WebSocket connection when component unmounts
    return () => {
      if (socket.current) {
        socket.current.close()
      }
    }
  }, []) // Empty dependency array ensures this only runs once when the component mounts

  // Handle sending a new message
  const handleSendMessage = () => {
    const newMessage = {
      username: "User1", // Use dynamic username if needed
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
