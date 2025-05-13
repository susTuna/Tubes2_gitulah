import type { NextApiRequest, NextApiResponse } from 'next'
import httpProxy from 'http-proxy'

export const config = {
  api: {
    bodyParser: false, // Needed for streaming the body directly
    externalResolver: true,
  },
}

const proxy = httpProxy.createProxyServer()

export default function handler(req: NextApiRequest, res: NextApiResponse) {
  return new Promise<void>((resolve, reject) => {
    proxy.web(req, res, {
      target: process.env.BACKEND_PUBLIC_API_URL,
      changeOrigin: true,
      ignorePath: false,
      secure: false,
    })

    proxy.once('proxyRes', () => resolve())
    proxy.once('error', (err) => reject(err))
    console.log("Proxying to:", `${process.env.BACKEND_PUBLIC_API_URL}${req.url?.replace('/api/proxy', '')}`)

  })
}

export const fetchFromBackend = async (path: string, options?: RequestInit) => {
    const res = await fetch(`/api/proxy${path}`, options)
    return res
  }