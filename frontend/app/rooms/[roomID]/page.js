"use client"

import { useEffect, useState, useRef, useContext } from "react"
import { useParams } from "next/navigation"
import axios from "axios"
import AuthContext from "../../context/AuthContext"

export default function RoomPage() {
  const params = useParams()
  const roomID = params.roomID
  const [messages, setMessages] = useState([])
  const [input, setInput] = useState("")
  const { isAuthenticated, loading } = useContext(AuthContext)
  const messageEndRef = useRef(null)
  const ws = useRef(null) // WebSocket reference

  // Fetch existing messages when component mounts
  useEffect(() => {
    const fetchMessages = async () => {
      try {
        const res = await axios.get(`/api/rooms/${roomID}/messages`, {
          withCredentials: true,
        })
        setMessages(res.data) // Assuming messages are ordered by newest first
      } catch (err) {
        console.error("Error fetching messages:", err)
      }
    }
    fetchMessages()
  }, [])

  // Establish WebSocket connection
  useEffect(() => {
    if (!isAuthenticated) return

    // Initialize WebSocket connection
    ws.current = new WebSocket(`ws://localhost:8080/ws/${roomID}`)

    ws.current.onopen = () => {
      console.log(`Connected to room ${roomID}`)
      // Optionally, send authentication details if required
      // Example: ws.current.send(JSON.stringify({ token: 'YOUR_JWT_TOKEN' }));
    }

    ws.current.onmessage = (event) => {
      try {
        const newMessage = JSON.parse(event.data)
        if (newMessage.error) {
          console.error("WebSocket Error:", newMessage.error)
          return
        }
        setMessages((prevMessages) => [...prevMessages, newMessage])
        scrollToBottom()
      } catch (err) {
        console.error("Error parsing WebSocket message:", err)
      }
    }

    ws.current.onclose = () => {
      console.log(`Disconnected from room ${roomID}`)
    }

    ws.current.onerror = (err) => {
      console.error("WebSocket error:", err)
    }

    // Cleanup on component unmount
    return () => {
      if (ws.current) {
        ws.current.close()
      }
    }
  }, [])

  // Scroll to the latest message
  const scrollToBottom = () => {
    messageEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }

  // Handle sending a message
  const handleSendMessage = async (e) => {
    e.preventDefault()
    if (input.trim() === "") return

    const messagePayload = { content: input }

    try {
      // Send message via HTTP POST to ensure it's saved
      await axios.post(`/api/rooms/${roomID}/messages`, messagePayload, {
        withCredentials: true,
      })

      // Optionally, send the message via WebSocket for immediate broadcast
      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        ws.current.send(JSON.stringify(messagePayload))
      }

      setInput("")
    } catch (err) {
      console.error("Failed to send message:", err)
    }
  }

  if (loading) {
    return <p>Loading...</p>
  }

  if (!isAuthenticated) {
    return <p>You must be logged in to view this room.</p>
  }

  return (
    <div className="min-h-screen p-8 bg-gray-100 flex flex-col">
      <h1 className="text-3xl font-bold mb-6">Room: {roomID}</h1>
      <div className="flex-1 overflow-y-auto bg-white p-4 rounded-lg shadow-md mb-4">
        {messages.map((message) => (
          <p key={message.id}>
            <strong>{message.username}:</strong> {message.content}
          </p>
        ))}
        <div ref={messageEndRef} />
      </div>
      <form onSubmit={handleSendMessage} className="flex space-x-4">
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Type a message..."
          className="flex-1 p-2 border border-gray-300 rounded"
          required
        />
        <button
          type="submit"
          className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
        >
          Send
        </button>
      </form>
    </div>
  )
}
