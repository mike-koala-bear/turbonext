"use client"

import Link from "next/link"

export default function Home() {
  return (
    <div>
      <h1>Welcome to the Chat App</h1>
      <Link href="/chat">Go to Chat Room</Link>
    </div>
  )
}
