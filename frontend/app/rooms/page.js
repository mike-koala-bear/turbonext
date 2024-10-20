"use client"

import { useState, useEffect } from "react"
import axios from "axios"
import Link from "next/link"

export default function RoomList() {
  const [rooms, setRooms] = useState([])

  useEffect(() => {
    const fetchRooms = async () => {
      const res = await axios.get("/api/rooms", { withCredentials: true })
      setRooms(res.data)
    }
    fetchRooms()
  }, [])

  return (
    <div className="m-10">
      <h2 className="text-3xl font-medium">Available Rooms</h2>

      <ul className="mt-4 list-disc">
        {rooms.map((room) => (
          <li key={room.id}>
            <Link
              href={`/rooms/${room.id}`}
              className="text-blue-500 hover:text-blue-400 hover:underline"
            >
              {room.name}
            </Link>
          </li>
        ))}
      </ul>
    </div>
  )
}
