import type { IRequest } from 'itty-router';

/**
 * 404 Not Found handler
 */
export async function handleNotFound(request: IRequest): Promise<Response> {
  const response = {
    success: false,
    error: 'Not Found',
    message: `The requested resource ${request.url} was not found`,
    timestamp: new Date().toISOString(),
  };

  return new Response(JSON.stringify(response, null, 2), {
    status: 404,
    headers: {
      'Content-Type': 'application/json',
    },
  });
}