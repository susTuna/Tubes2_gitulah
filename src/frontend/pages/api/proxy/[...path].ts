import type { NextApiRequest, NextApiResponse } from 'next'
import { createProxyMiddleware } from 'http-proxy-middleware'

export const config = {
  api: {
    bodyParser: false, // Required for streaming body
    externalResolver: true,
  },
}

const proxy = createProxyMiddleware({
  target: process.env.BACKEND_PUBLIC_API_URL,
  changeOrigin: true,
  secure: false,
  pathRewrite: { '^/api/proxy': '' }, // This will strip '/api/proxy' from the request URL
})

export default function handler(req: NextApiRequest, res: NextApiResponse) {
  if (!process.env.BACKEND_PUBLIC_API_URL) {
    return res.status(500).json({ error: 'Backend URL is not configured' })
  }

  console.log("Proxying to:", `${process.env.BACKEND_PUBLIC_API_URL}${req.url?.replace('/api/proxy', '')}`)

  // Call the proxy middleware to forward the request
  return proxy(req, res)
}

export const fetchFromBackend = async (path: string, options?: RequestInit) => {
    const res = await fetch(`/api/proxy${path}`, options)
    return res
  }