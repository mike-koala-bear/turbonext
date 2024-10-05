/** @type {import('next').NextConfig} */
const nextConfig = {
  // Proxying API requests during development
  reactStrictMode: false,
  async rewrites() {
    return [
      {
        source: "/api/:path*",
        destination: "http://localhost:8080/:path*", // Proxy to Go backend
      },
    ]
  },
}

export default nextConfig
