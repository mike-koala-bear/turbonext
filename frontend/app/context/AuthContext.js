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
            return status === 200 || status === 401
          },
        })

        if (res.status === 200) {
          setIsAuthenticated(true)
        } else {
          setIsAuthenticated(false)
        }
      } catch (error) {
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
