"use client"

import { createContext, useState, useEffect } from "react"
import axios from "axios"

const AuthContext = createContext()

export function AuthProvider({
  children,
  isAuthenticated: initialAuth = false,
}) {
  const [isAuthenticated, setIsAuthenticated] = useState(initialAuth)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const validateToken = async () => {
      try {
        const res = await axios.get("/api/check-auth", {
          withCredentials: true,
          validateStatus: function (status) {
            // Accept status codes 200 and 401
            return status === 200 || status === 401
          },
        })

        if (res.status === 200) {
          // User is authenticated
          setIsAuthenticated(true)
        } else if (res.status === 401) {
          // User is not authenticated
          setIsAuthenticated(false)
        }
      } catch (error) {
        console.error(
          "An error occurred while checking authentication status:",
          error
        )
        setIsAuthenticated(false)
      } finally {
        setLoading(false)
      }
    }

    validateToken()
  }, [])

  return (
    <AuthContext.Provider
      value={{ isAuthenticated, setIsAuthenticated, loading }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export default AuthContext
