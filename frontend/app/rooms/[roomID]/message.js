import { useState, useEffect, useRef } from "react"
import { getRelativeTime } from "../../utils"

export default function Message({ message, isCurrentUser }) {
  const [relativeTime, setRelativeTime] = useState("")
  const timerRef = useRef(null)

  const updateRelativeTime = () => {
    const { relativeTime, nextInterval } = getRelativeTime(message.created_at)
    setRelativeTime(relativeTime)

    // Clear any existing timers
    if (timerRef.current) {
      clearTimeout(timerRef.current)
    }

    // Schedule the next update
    timerRef.current = setTimeout(updateRelativeTime, nextInterval * 1000)
  }

  useEffect(() => {
    updateRelativeTime()

    // Clean up on unmount
    return () => {
      if (timerRef.current) {
        clearTimeout(timerRef.current)
      }
    }
  }, [message.created_at])

  return (
    <div className={`flex ${isCurrentUser ? "justify-end" : "justify-start"}`}>
      <div
        className={`relative max-w-xs rounded-lg p-3 shadow-lg ${
          isCurrentUser ? "bg-blue-500 text-white" : "bg-gray-200 text-gray-800"
        }`}
      >
        {/* Tail */}
        <div
          className={`absolute bottom-0 w-0 h-0 border-t-8 border-t-transparent ${
            isCurrentUser
              ? "right-0 transform translate-x-1/2 border-l-8 border-l-blue-500"
              : "left-0 transform -translate-x-1/2 border-r-8 border-r-gray-200"
          }`}
        ></div>

        {/* Header (Username and Time) */}
        <div className="flex items-center mb-1">
          <span className="text-md font-semibold">{message.username}</span>
          <span
            className={`text-xs ${
              isCurrentUser ? "text-gray-300" : "text-gray-800"
            } ml-2`}
          >
            {relativeTime}
          </span>
        </div>

        {/* Message Content */}
        <div className="break-words leading-relaxed">{message.content}</div>
      </div>
    </div>
  )
}
