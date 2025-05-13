import type { NextApiRequest, NextApiResponse } from 'next'
import { createProxyMiddleware } from 'http-proxy-middleware'

export const config = {
  api: {
    bodyParser: false, // Required for streaming body
    externalResolver: true,
  },
}

const proxy = createProxyMiddleware({
  target: process.env.NEXT_PUBLIC_BACKEND_PUBLIC_API_URL,
  changeOrigin: true,
  secure: false,
  pathRewrite: { '^/api/proxy': '' },
})

export default function handler(req: NextApiRequest, res: NextApiResponse) {
  if (!process.env.NEXT_PUBLIC_BACKEND_PUBLIC_API_URL) {
    return res.status(500).json({ error: 'Backend URL is not configured' })
  }

  console.log("Proxying to:", `${process.env.NEXT_PUBLIC_BACKEND_PUBLIC_API_URL}${req.url?.replace('/api/proxy', '')}`)

  // Call the proxy middleware to forward the request
  return proxy(req, res)
}

export const fetchFromBackend = async (path: string, options?: RequestInit) => {
    const res = await fetch(`/api/proxy${path}`, options)
    return res
  }