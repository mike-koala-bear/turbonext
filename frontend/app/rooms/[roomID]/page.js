"use client"

import { useEffect, useState, useRef, useContext } from "react"
import { useParams } from "next/navigation"
import axios from "axios"
import AuthContext from "../../context/AuthContext"
import Message from "./message"

export default function RoomPage() {
  const params = useParams()
  const roomID = params.roomID
  const [roomName, setRoomName] = useState("")
  const [messages, setMessages] = useState([])
  const [input, setInput] = useState("")
  const { isAuthenticated, loading, username } = useContext(AuthContext)
  const messageEndRef = useRef(null)
  const ws = useRef(null)

  // Fetch existing messages when component mounts
  useEffect(() => {
    const fetchMessages = async () => {
      try {
        const res = await axios.get(`/api/rooms/${roomID}/messages`, {
          withCredentials: true,
        })
        console.log("API Response:", res.data)
        setMessages(res.data.messages)
        setRoomName(res.data.room_name)
      } catch (err) {
        console.error("Error fetching messages:", err)
      }
    }

    fetchMessages()
  }, [roomID])

  // Scroll to the last message when messages are fetched or updated
  useEffect(() => {
    scrollToBottom()
  }, [messages])

  // Establish WebSocket connection
  useEffect(() => {
    if (!isAuthenticated) return

    const setupWebSocket = () => {
      ws.current = new WebSocket(`ws://localhost:8080/ws/${roomID}`)

      ws.current.onopen = () => {
        console.log(`Connected to room ${roomID}`)
      }

      ws.current.onmessage = (event) => {
        try {
          const newMessage = JSON.parse(event.data)
          if (newMessage.error) {
            console.error("WebSocket Error:", newMessage.error)
            return
          }
          setMessages((prevMessages) => [...prevMessages, newMessage]) // Append new messages at the end
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
    }

    setupWebSocket()

    return () => {
      if (ws.current) {
        ws.current.close()
      }
    }
  }, [isAuthenticated, roomID])

  // Scroll to the last message
  const scrollToBottom = () => {
    messageEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }

  // Handle sending a message via HTTP POST
  const handleSendMessage = async (e) => {
    e.preventDefault()
    if (input.trim() === "") return

    const messagePayload = { content: input }

    try {
      const res = await axios.post(
        `/api/rooms/${roomID}/messages`,
        messagePayload,
        {
          withCredentials: true,
        }
      )

      if (res.status === 200) {
        console.log("Message sent successfully")
        setInput("")
      } else {
        console.error("Failed to send message:", res.data.error)
      }
    } catch (err) {
      console.error("Failed to send message via HTTP POST:", err)
    }
  }

  if (loading) {
    return <p>Loading...</p>
  }

  if (!isAuthenticated) {
    return <p>You must be logged in to view this room.</p>
  }

  return (
    <div className="min-h-screen flex flex-col bg-gray-100">
      {/* Header */}
      <div className="bg-white shadow-md p-4 sticky top-0 z-10">
        <h1 className="text-2xl font-bold text-center">Room: {roomName}</h1>
      </div>

      {/* Messages List */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4 pb-24">
        {messages.map((message) => {
          const messageUsername = message.username?.trim().toLowerCase() || ""
          const currentUsername = username?.trim().toLowerCase() || ""
          const isCurrentUser = messageUsername === currentUsername

          return (
            <Message
              key={message.id}
              message={message}
              isCurrentUser={isCurrentUser}
            />
          )
        })}
        {/* Spacer at the bottom */}
        <div ref={messageEndRef} />
      </div>

      <form
        onSubmit={handleSendMessage}
        className="fixed bottom-0 left-0 w-full flex items-center space-x-4 p-4 bg-white border-t border-gray-300"
      >
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Type a message..."
          className="flex-1 p-3 border border-gray-300 rounded-full focus:outline-none focus:ring-2 focus:ring-blue-400"
          required
        />
        <button
          type="submit"
          className="bg-blue-500 text-white px-4 py-2 rounded-full hover:bg-blue-600"
        >
          Send
        </button>
      </form>
    </div>
  )
}
