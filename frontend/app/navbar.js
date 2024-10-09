// app/navbar.js
"use client"

import Link from "next/link"
import { useRouter } from "next/navigation"
import { useContext } from "react"
import axios from "axios"
import AuthContext from "./context/AuthContext"

export default function Navbar() {
  const router = useRouter()
  const { isAuthenticated, setIsAuthenticated } = useContext(AuthContext)

  const handleLogout = () => {
    axios
      .post("/api/logout", {}, { withCredentials: true })
      .then(() => {
        setIsAuthenticated(false)
        router.push("/login")
      })
      .catch((error) => {
        console.error("An error occurred during logout:", error)
      })
  }

  return (
    <nav className="bg-gray-800 p-4">
      <div className="max-w-7xl mx-auto flex justify-between">
        <Link href="/" className="text-white">
          Home
        </Link>
        <div>
          {isAuthenticated && (
            <>
              <Link href="/create-room" className="text-white mr-4">
                Create Room
              </Link>
              <Link href="/rooms" className="text-white mr-4">
                Rooms
              </Link>
            </>
          )}
          {!isAuthenticated ? (
            <>
              <Link href="/login" className="text-white mr-4">
                Login
              </Link>
              <Link href="/signup" className="text-white">
                Signup
              </Link>
            </>
          ) : (
            <button onClick={handleLogout} className="text-white">
              Logout
            </button>
          )}
        </div>
      </div>
    </nav>
  )
}
