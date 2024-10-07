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
    <div>
      <h2>Available Rooms</h2>
      <ul>
        {rooms.map((room) => (
          <li key={room.id}>
            <Link href={`/rooms/${room.id}`}>{room.name}</Link>
          </li>
        ))}
      </ul>
    </div>
  )
}
