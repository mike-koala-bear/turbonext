"use client"

import { AuthProvider } from "./context/AuthContext"

export default function Providers({ children, isAuthenticated }) {
  return (
    <AuthProvider isAuthenticated={isAuthenticated}>{children}</AuthProvider>
  )
}
