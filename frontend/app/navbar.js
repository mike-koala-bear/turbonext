"use client"

import Link from "next/link"
import { useRouter } from "next/navigation"
import { useEffect, useState } from "react"
import axios from "axios"

export default function Navbar() {
  const router = useRouter()
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [hasMounted, setHasMounted] = useState(false)

  useEffect(() => {
    // Check authentication status by making an API request
    const checkAuthStatus = async () => {
      try {
        // Send a request to a protected route
        const res = await axios.get("/api/check-auth", {
          withCredentials: true,
        }) // Token sent via cookie
        if (res.status === 200) {
          setIsAuthenticated(true) // User is authenticated
        }
      } catch (error) {
        setIsAuthenticated(false) // User is not authenticated
      }
    }

    checkAuthStatus()
    setHasMounted(true) // Set the component as mounted
  }, [])

  // Prevent rendering until after the component has mounted
  if (!hasMounted) {
    return null // or a loading indicator if you prefer
  }

  const handleLogout = () => {
    // Logout logic (e.g., clear cookies, or backend endpoint for logout)
    axios.post("/api/logout", {}, { withCredentials: true }).then(() => {
      router.push("/login")
    })
  }

  return (
    <nav className="bg-gray-800 p-4">
      <div className="max-w-7xl mx-auto flex justify-between">
        <Link href="/" className="text-white">
          Home
        </Link>
        <div>
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
