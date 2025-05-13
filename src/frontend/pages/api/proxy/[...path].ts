import type { NextApiRequest, NextApiResponse } from 'next'

// Configure API to handle request bodies
export const config = {
  api: {
    bodyParser: true,
  },
}

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  const backendUrl = process.env.NEXT_PUBLIC_BACKEND_PUBLIC_API_URL;
  
  if (!backendUrl) {
    return res.status(500).json({ error: 'Backend URL is not configured' });
  }

  // Get the path from the request
  const pathSegments = req.query.path as string[];
  let path = '/' + pathSegments.join('/');
  
  // Ensure fullrecipe endpoint has trailing slash
  if (path === '/fullrecipe') {
    path = '/fullrecipe/';
  }

  console.log(`${req.method} request - Proxying to: ${backendUrl}${path}`);
  
  try {
    // Build fetch options based on the original request
    const fetchOptions: RequestInit = {
      method: req.method,
      headers: {
        'Content-Type': req.headers['content-type'] || 'application/json',
      },
    };

    // Add body for POST, PUT, PATCH requests
    if (['POST', 'PUT', 'PATCH'].includes(req.method || '')) {
      fetchOptions.body = JSON.stringify(req.body);
    }

    // Forward the request to the backend
    const response = await fetch(`${backendUrl}${path}`, fetchOptions);
    
    // Get response data
    const data = await response.text();
    const contentType = response.headers.get('content-type') || 'application/json';
    
    // Set status and headers
    res.status(response.status);
    res.setHeader('Content-Type', contentType);
    
    // Return the response
    if (contentType.includes('application/json')) {
      try {
        return res.send(JSON.parse(data));
      } catch {
        return res.send(data);
      }
    } else {
      return res.send(data);
    }
    
  } catch (error) {
    console.error('Proxy error:', error);
    return res.status(500).json({ 
      error: 'Failed to proxy request',
      details: error instanceof Error ? error.message : 'Unknown error' 
    });
  }
}

// Helper function for frontend components to use
export const fetchFromBackend = async (path: string, options?: RequestInit) => {
  // Ensure fullrecipe endpoint has trailing slash
  if (path === '/fullrecipe') {
    path = '/fullrecipe/';
  }
  return fetch(`/api/proxy${path}`, options);
}