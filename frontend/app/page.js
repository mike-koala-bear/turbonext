"use client"

import Link from "next/link"
import { useContext } from "react"
import AuthContext from "./context/AuthContext"

export default function Home() {
  const { isAuthenticated } = useContext(AuthContext)

  return (
    <div>
      <h1>Welcome to the Chat App</h1>

      {isAuthenticated && <p>You are logged in!</p>}
    </div>
  )
}
