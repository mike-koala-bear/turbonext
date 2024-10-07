"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import axios from "axios"

export default function CreateRoom() {
  const [roomName, setRoomName] = useState("")
  const [error, setError] = useState("")
  const router = useRouter()

  const handleCreateRoom = async () => {
    try {
      const res = await axios.post(
        "/api/rooms",
        { name: roomName },
        { withCredentials: true }
      )
      router.push(`/rooms/${res.data.id}`)
    } catch (err) {
      setError("Failed to create room")
    }
  }

  return (
    <div>
      <h2>Create a New Room</h2>
      <input
        type="text"
        value={roomName}
        onChange={(e) => setRoomName(e.target.value)}
        placeholder="Enter room name"
      />
      <button onClick={handleCreateRoom}>Create Room</button>
      {error && <p>{error}</p>}
    </div>
  )
}
