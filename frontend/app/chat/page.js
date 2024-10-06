"use client"

import { useEffect, useState, useRef } from "react"
import axios from "axios"
import { useRouter } from "next/navigation"

export default function ChatPage() {
  const [messages, setMessages] = useState([]) // Store chat messages
  const [input, setInput] = useState("") // Store user's input
  const socket = useRef(null) // WebSocket reference
  const router = useRouter()
  const reconnectDelay = useRef(5000) // Delay for reconnecting WebSocket

  // Fetch chat messages and establish WebSocket connection
  useEffect(() => {
    // Fetch persisted chat messages from the server
    axios
      .get("/api/messages", { withCredentials: true })
      .then((response) => {
        setMessages(response.data) // Store fetched messages in state
      })
      .catch((error) => {
        console.error("Error fetching messages:", error)
        router.push("/login") // Redirect to login on error
      })

    // WebSocket connection and reconnection logic
    const connectWebSocket = () => {
      socket.current = new WebSocket("ws://localhost:8080/ws")

      // On WebSocket connection, send JWT token for authentication
      socket.current.onopen = async () => {
        console.log("WebSocket connected")

        // Get the JWT token from the server for WebSocket authentication
        const tokenResponse = await axios.get("/api/get-token", {
          withCredentials: true,
        })
        const token = tokenResponse.data.token

        // Send the token as the first WebSocket message for authentication
        socket.current.send(
          JSON.stringify({
            token: token,
          })
        )
      }

      // Handle incoming WebSocket messages
      socket.current.onmessage = (event) => {
        console.log("WebSocket received a message:", event.data)
        const newMessage = JSON.parse(event.data)

        // Update chat messages state
        setMessages((prevMessages) => [newMessage, ...prevMessages])
      }

      // Handle WebSocket closure and reconnect
      socket.current.onclose = (event) => {
        console.error("WebSocket closed:", event)
        if (!event.wasClean) {
          console.log(
            `Reconnecting WebSocket in ${reconnectDelay.current / 1000} seconds`
          )
          setTimeout(connectWebSocket, reconnectDelay.current)
        }
      }

      // Handle WebSocket errors
      socket.current.onerror = (error) => {
        console.error("WebSocket error:", error)
        socket.current.close() // Close WebSocket on error
      }
    }

    connectWebSocket() // Start WebSocket connection

    // Clean up WebSocket connection when the component is unmounted
    return () => {
      if (socket.current) {
        console.log("Closing WebSocket")
        socket.current.close()
      }
    }
  }, [])

  // Send a new message via WebSocket
  const handleSendMessage = () => {
    const newMessage = {
      content: input, // Send the user's message input
    }
    console.log("sending message: ", newMessage)
    // Check if WebSocket is open before sending the message
    if (socket.current && socket.current.readyState === WebSocket.OPEN) {
      console.log("socket is open and is sending message")
      socket.current.send(JSON.stringify(newMessage))
    }
    console.log("clearing input")
    setInput("") // Clear input field after sending message
  }

  return (
    <div className="min-h-screen bg-gray-100 p-8 flex flex-col items-center">
      <h1 className="text-3xl font-bold mb-4">Chat Room</h1>

      <div className="bg-white shadow-lg rounded-lg p-6 w-full max-w-lg">
        {/* Display chat messages */}
        <div className="mb-4 h-64 overflow-y-auto border p-4 bg-gray-50 rounded">
          {messages.map((msg, index) => (
            <div key={index} className="mb-2">
              <strong>{msg.username}:</strong> {msg.content}
              <br />
              <small className="text-gray-500">
                {new Date(msg.created_at).toLocaleString()}
              </small>
            </div>
          ))}
        </div>

        {/* Input for new messages */}
        <div className="flex">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Type a message..."
            className="flex-grow p-2 border border-gray-300 rounded-l-md focus:outline-none"
          />
          <button
            onClick={handleSendMessage}
            className="bg-blue-500 text-white px-4 py-2 rounded-r-md hover:bg-blue-600"
          >
            Send
          </button>
        </div>
      </div>
    </div>
  )
}
