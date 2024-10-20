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
    <div className="m-10">
      <h2 className="text-3xl font-medium">Create a New Room</h2>
      <div className="mt-6">
        <input
          type="text"
          value={roomName}
          onChange={(e) => setRoomName(e.target.value)}
          placeholder="Enter room name"
          className="border border-gray-300 p-2 px-4 rounded-md"
        />
        <button
          onClick={handleCreateRoom}
          className="ml-4 inline-flex items-center justify-center rounded-xl bg-green-600
           py-3 px-4 font-dm text-base font-medium text-white shadow-lg shadow-green-400/75 
           transition-transform duration-200 ease-in-out hover:scale-[1.05]"
        >
          Create Room
        </button>
        {error && <p>{error}</p>}
      </div>
    </div>
  )
}
