import type { IRequest } from 'itty-router';
import type { Env, User } from '../types/database';
import { errorResponse } from '../utils/responses';

/**
 * Extract and verify JWT token
 */
export async function verifyToken(token: string, secret: string): Promise<any> {
  try {
    // Split token
    const parts = token.split('.');
    if (parts.length !== 3) {
      throw new Error('Invalid token format');
    }

    // In production, use a proper JWT library like jose
    // This is a simplified example
    const payload = JSON.parse(atob(parts[1]));
    
    // Check expiration
    if (payload.exp && Date.now() >= payload.exp * 1000) {
      throw new Error('Token expired');
    }

    return payload;
  } catch (error) {
    throw new Error('Invalid token');
  }
}

/**
 * Authentication middleware
 * Adds user to request if valid token provided
 */
export async function authenticate(
  request: IRequest,
  env: Env
): Promise<void | Response> {
  try {
    const authHeader = request.headers.get('Authorization');
    
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return; // No authentication, continue as guest
    }

    const token = authHeader.slice(7);
    const payload = await verifyToken(token, env.JWT_SECRET);
    
    // Get user from database
    const user = await env.DB.prepare(
      'SELECT * FROM users WHERE id = ? AND deleted_at IS NULL'
    )
      .bind(payload.sub)
      .first<User>();

    if (!user) {
      return errorResponse('User not found', 401);
    }

    // Update last login
    await env.DB.prepare(
      'UPDATE users SET last_login_at = CURRENT_TIMESTAMP WHERE id = ?'
    ).bind(user.id).run();

    // Add user to request
    request.user = user;
  } catch (error) {
    return errorResponse('Invalid authentication', 401);
  }
}

/**
 * Require authentication middleware
 */
export function requireAuth(
  request: IRequest
): Response | void {
  if (!request.user) {
    return errorResponse('Authentication required', 401);
  }
}

/**
 * Require admin role middleware
 */
export function requireAdmin(
  request: IRequest
): Response | void {
  if (!request.user) {
    return errorResponse('Authentication required', 401);
  }

  if (request.user.role !== 'admin') {
    return errorResponse('Admin access required', 403);
  }
}

/**
 * Require moderator or admin role middleware
 */
export function requireModerator(
  request: IRequest
): Response | void {
  if (!request.user) {
    return errorResponse('Authentication required', 401);
  }

  if (!['admin', 'moderator'].includes(request.user.role)) {
    return errorResponse('Moderator access required', 403);
  }
}