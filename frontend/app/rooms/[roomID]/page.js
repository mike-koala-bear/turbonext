"use client"

import { useEffect, useState, useRef, useContext } from "react"
import { useParams } from "next/navigation"
import axios from "axios"
import AuthContext from "../../context/AuthContext"
import { format } from "date-fns"

export default function RoomPage() {
  const params = useParams()
  const roomID = params.roomID
  const [roomName, setRoomName] = useState("")
  const [messages, setMessages] = useState([])
  const [input, setInput] = useState("")
  const { isAuthenticated, loading } = useContext(AuthContext)
  const messageStartRef = useRef(null)
  const ws = useRef(null)

  // Fetch existing messages when component mounts
  useEffect(() => {
    const fetchMessages = async () => {
      try {
        const res = await axios.get(`/api/rooms/${roomID}/messages`, {
          withCredentials: true,
        })
        console.log("API Response:", res.data)
        setMessages(res.data.messages.reverse()) // Keep reverse to maintain order
        setRoomName(res.data.room_name)
        scrollToTop()
      } catch (err) {
        console.error("Error fetching messages:", err)
      }
    }

    fetchMessages()
  }, [roomID])

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
          setMessages((prevMessages) => [newMessage, ...prevMessages])
          // Scroll to top since new message is at the top
          scrollToTop()
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

  // Scroll to the latest message (now at the top)
  const scrollToTop = () => {
    messageStartRef.current?.scrollIntoView({ behavior: "smooth" })
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
    <div className="min-h-screen p-8 bg-gray-100 flex flex-col">
      <h1 className="text-3xl font-bold mb-6">Room: {roomName}</h1>
      <div className="flex-1 overflow-y-auto bg-white p-4 rounded-lg shadow-md mb-4">
        <div ref={messageStartRef} />
        {messages.map((message) => (
          <div key={message.id} className="mb-4">
            <p className="text-sm text-gray-500">
              <strong>{message.username}</strong> •{" "}
              {format(new Date(message.created_at), "PPPpp")}
            </p>
            <p className="text-base">{message.content}</p>
          </div>
        ))}
      </div>
      <form onSubmit={handleSendMessage} className="flex space-x-4 mt-4">
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
