import {
  differenceInSeconds,
  differenceInMinutes,
  differenceInHours,
  differenceInDays,
  differenceInMonths,
  differenceInYears,
} from "date-fns"

export function getRelativeTime(timestamp) {
  const now = new Date()
  const time = new Date(timestamp)
  const secondsDifference = differenceInSeconds(now, time)
  let relativeTime = ""
  let nextInterval = 0

  if (secondsDifference < 60) {
    // Less than 1 minute ago
    relativeTime = "Just now"
    nextInterval = 60 - secondsDifference
  } else if (secondsDifference < 3600) {
    // Between 1 minute and 1 hour
    const minutesDifference = differenceInMinutes(now, time)
    relativeTime = `${minutesDifference} minute${
      minutesDifference !== 1 ? "s" : ""
    } ago`
    nextInterval = 60 - (secondsDifference % 60)
  } else if (secondsDifference < 86400) {
    // Between 1 hour and 24 hours
    const hoursDifference = differenceInHours(now, time)
    relativeTime = `${hoursDifference} hour${
      hoursDifference !== 1 ? "s" : ""
    } ago`
    nextInterval = 3600 - (secondsDifference % 3600)
  } else if (secondsDifference < 2592000) {
    // Between 1 day and 30 days
    const daysDifference = differenceInDays(now, time)
    relativeTime = `${daysDifference} day${daysDifference !== 1 ? "s" : ""} ago`
    nextInterval = 86400 - (secondsDifference % 86400)
  } else if (secondsDifference < 31536000) {
    // Between 1 month and 1 year
    const monthsDifference = differenceInMonths(now, time)
    relativeTime = `${monthsDifference} month${
      monthsDifference !== 1 ? "s" : ""
    } ago`
    nextInterval = 2592000 - (secondsDifference % 2592000)
  } else {
    // More than 1 year
    const yearsDifference = differenceInYears(now, time)
    relativeTime = `${yearsDifference} year${
      yearsDifference !== 1 ? "s" : ""
    } ago`
    nextInterval = 31536000 - (secondsDifference % 31536000)
  }

  // Ensure nextInterval is at least 1 second
  nextInterval = Math.max(nextInterval, 1)

  return { relativeTime, nextInterval }
}
