/**
 * NawthTech Cloudflare Worker
 * Main entry point with D1 Database integration
 */

import { Router } from 'itty-router';
import { error, json, missing, status } from 'itty-router-extras';
import type { IRequest } from 'itty-router';

// Import handlers
import { handleHealthCheck } from './handlers/health';
import { handleCORS } from './middleware/cors';
import { authenticate, requireAuth, requireAdmin } from './middleware/auth';
import { validateRequest } from './middleware/validation';
import {
  registerUser,
  loginUser,
  getCurrentUser,
  updateUser,
  listUsers,
  getUserById,
} from './handlers/api/v1/auth';
import {
  createService,
  getServices,
  getServiceById,
  updateService,
  deleteService,
} from './handlers/api/v1/services';
import {
  generateAI,
  getAIQuota,
  getAIRequests,
} from './handlers/api/v1/ai';
import {
  handleEmailWebhook,
  getEmailLogs,
  getEmailConfig,
} from './handlers/webhooks/email';
import { handleStaticFile } from './handlers/static';
import { handleNotFound } from './handlers/notfound';
import type { Env, User } from './types/database';
import type { APIResponse } from './utils/responses';
import * as Sentry from '@sentry/cloudflare';
import { initSentry, withSentryErrorBoundary, captureMessage } from './sentry-setup.js';

export interface Env {
  // Your bindings
  DB: D1Database;
  AI: Ai;
  // Sentry
  SENTRY_DSN?: string;
  SENTRY_ENVIRONMENT?: string;
}

// Initialize Sentry with environment variables
const SentryInstance = Sentry.withSentry((env: Env) => ({
  dsn: env.SENTRY_DSN || "https://703dc8c9404510702c2c20ce3aba24d4@o4510508331892736.ingest.de.sentry.io/4510508452413520",
  environment: env.SENTRY_ENVIRONMENT || 'production',
  release: 'nawthtech-worker@1.0.0',
  sendDefaultPii: true,
  tracesSampleRate: 0.2,
  integrations: [
    new Sentry.CloudflareIntegration(),
  ],
  beforeSend(event) {
    // Add custom context for nawthtech
    event.tags = {
      ...event.tags,
      platform: 'cloudflare-workers',
      service: 'nawthtech-social-growth',
      app: 'nawthtech',
      component: 'main-worker',
    };
    
    // Add user context if available
    if (event.request?.headers?.get('x-user-id')) {
      event.user = {
        id: event.request.headers.get('x-user-id'),
        ip_address: '{{auto}}',
      };
    }
    
    return event;
  },
}));

export default SentryInstance({
  async fetch(request: Request, env: Env, ctx: ExecutionContext) {
    // Start Sentry transaction
    const transaction = Sentry.startTransaction({
      name: `${request.method} ${new URL(request.url).pathname}`,
      op: 'http.server',
    });
    
    Sentry.configureScope(scope => {
      scope.setSpan(transaction);
      scope.setTag('http.method', request.method);
      scope.setTag('http.url', request.url);
      scope.setTag('app.feature', 'social-growth-platform');
    });
    
    try {
      // Test error (optional - for verification)
      if (new URL(request.url).pathname === '/test-error') {
        setTimeout(() => {
          throw new Error('Test error for Sentry verification - nawthtech social growth platform');
        }, 100);
        return new Response('Test error triggered');
      }
      
      // Health check endpoint
      if (new URL(request.url).pathname === '/health') {
        captureMessage('Health check performed', 'info', {
          timestamp: new Date().toISOString(),
        });
        return new Response('OK', { status: 200 });
      }
      
      // Your main logic here
      const response = await handleRequest(request, env, ctx);
      
      // Capture successful request
      captureMessage('Request processed successfully', 'info', {
        path: new URL(request.url).pathname,
        method: request.method,
        status: response.status,
      });
      
      return response;
      
    } catch (error) {
      // Capture error with nawthtech context
      Sentry.captureException(error, {
        tags: {
          error_type: error instanceof Error ? error.constructor.name : 'Unknown',
          app_section: 'social-growth-worker',
          user_action: 'fetch_request',
        },
        extra: {
          request_url: request.url,
          request_method: request.method,
          timestamp: new Date().toISOString(),
          nawthtech_context: 'social-media-intelligence-platform',
        },
      });
      
      return new Response('Internal Server Error', { status: 500 });
    } finally {
      transaction.finish();
    }
  },
  
  // Scheduled handler (if you have cron triggers)
  async scheduled(event: ScheduledEvent, env: Env, ctx: ExecutionContext) {
    captureMessage('Scheduled task executed', 'info', {
      cron: event.cron,
      scheduledTime: event.scheduledTime,
      task: 'nawthtech-background-job',
    });
    
    try {
      // Your scheduled task logic
      await processBackgroundTasks(env);
    } catch (error) {
      Sentry.captureException(error, {
        tags: {
          task_type: 'scheduled',
          app_feature: 'background-processing',
        },
      });
      throw error;
    }
  },
} satisfies ExportedHandler<Env>);

// Your request handler with error boundary
const handleRequest = withSentryErrorBoundary(
  async (request: Request, env: Env, ctx: ExecutionContext) => {
    // Your existing handler logic
    return new Response('Hello from nawthtech social growth platform!');
  },
  'handleRequest'
);

// Background task processor
async function processBackgroundTasks(env: Env) {
  // Your background task logic
  console.log('Processing nawthtech background tasks...');
}
},
// Create router
const router = Router<IRequest, [Env, ExecutionContext]>();

// ============ Middleware ============

// Global CORS middleware
router.all('*', handleCORS);

// Global authentication middleware (adds user to request if token exists)
router.all('/api/*', authenticate);

// ============ Routes ============

// Health check
router.get('/health', handleHealthCheck);
router.get('/api/health', handleHealthCheck);

// Static files
router.get('/assets/*', handleStaticFile);
router.get('/favicon.ico', handleStaticFile);

// API v1 Routes
const apiV1 = router.basePath('/api/v1');

// Auth routes
apiV1.post('/auth/register', validateRequest('register'), registerUser);
apiV1.post('/auth/login', validateRequest('login'), loginUser);
apiV1.get('/auth/me', requireAuth, getCurrentUser);
apiV1.put('/auth/profile', requireAuth, validateRequest('updateProfile'), updateUser);

// Users routes (admin only)
apiV1.get('/users', requireAdmin, listUsers);
apiV1.get('/users/:id', requireAuth, getUserById);
apiV1.put('/users/:id', requireAdmin, validateRequest('updateUser'), updateUser);

// Services routes
apiV1.get('/services', requireAuth, getServices);
apiV1.post('/services', requireAuth, validateRequest('createService'), createService);
apiV1.get('/services/:id', requireAuth, getServiceById);
apiV1.put('/services/:id', requireAuth, validateRequest('updateService'), updateService);
apiV1.delete('/services/:id', requireAuth, deleteService);

// AI routes
apiV1.post('/ai/generate', requireAuth, validateRequest('generateAI'), generateAI);
apiV1.get('/ai/quota', requireAuth, getAIQuota);
apiV1.get('/ai/requests', requireAuth, getAIRequests);

// Email routes (admin only)
apiV1.get('/email/logs', requireAdmin, getEmailLogs);
apiV1.get('/email/config', requireAdmin, getEmailConfig);

// Webhooks
router.post('/webhooks/email', handleEmailWebhook);
router.post('/webhooks/stripe', async (request, env) => {
  // Stripe webhook handler
  return json({ received: true });
});

// Catch-all for SPA (if serving frontend)
router.get('*', handleStaticFile);

// ============ Error Handling ============

// 404 Not Found
router.all('*', handleNotFound);

// Global error handler
const handleError = (error: any): Response => {
  console.error('Unhandled error:', error);

  if (error.status) {
    return json(
      {
        success: false,
        error: error.message || 'An error occurred',
        code: error.code,
      },
      { status: error.status }
    );
  }

  // Internal server error
  return json(
    {
      success: false,
      error: 'Internal server error',
      message: env.ENVIRONMENT === 'development' ? error.message : undefined,
    },
    { status: 500 }
  );
};

// ============ Worker Entry Points ============

// Main fetch handler
export default {
  async fetch(
    request: Request,
    env: Env,
    ctx: ExecutionContext
  ): Promise<Response> {
    try {
      // Initialize database on first request (optional)
      ctx.waitUntil(initializeDatabase(env));

      // Handle request
      const response = await router.handle(request, env, ctx);
      return response;
    } catch (err) {
      return handleError(err);
    }
  },

  // Email Worker handler
  async email(message: ForwardableEmailMessage, env: Env, ctx: ExecutionContext): Promise<void> {
    try {
      console.log(`📧 Email received from: ${message.from}`);
      console.log(`📨 Subject: ${message.headers.get('subject')}`);

      // Parse configuration
      const allowList = env.EMAIL_ALLOWED_LIST 
        ? env.EMAIL_ALLOWED_LIST.split(',').map(e => e.trim().toLowerCase())
        : ['admin@nawthtech.com', 'support@nawthtech.com'];
      
      const forwardTo = env.EMAIL_FORWARD_TO || 'admin@nawthtech.com';
      const domain = 'nawthtech.com';

      // Check sender
      const senderEmail = message.from.toLowerCase().trim();
      let isAllowed = allowList.includes(senderEmail) || senderEmail.endsWith(`@${domain}`);

      if (!isAllowed) {
        console.log(`❌ Rejected email from: ${message.from}`);
        message.setReject('Address not allowed');
        return;
      }

      // Log to database
      try {
        await env.DB.prepare(`
          INSERT INTO email_logs (id, from_email, to_email, subject, status)
          VALUES (?, ?, ?, ?, ?)
        `)
          .bind(
            crypto.randomUUID(),
            message.from,
            Array.from(message.to).join(','),
            message.headers.get('subject') || 'No subject',
            'forwarded'
          )
          .run();
      } catch (error) {
        console.error('Failed to log email:', error);
      }

      // Forward email
      console.log(`✅ Forwarding email from ${message.from} to ${forwardTo}`);
      await message.forward(forwardTo);

      // Send acknowledgment (optional)
      if (env.ENVIRONMENT === 'production') {
        try {
          await message.reply(
            `Thank you for your email to NawthTech. Your message has been received and will be reviewed shortly.\n\nBest regards,\nNawthTech Team`
          );
        } catch (e) {
          console.log('Could not send acknowledgment:', e);
        }
      }

    } catch (error) {
      console.error('❌ Email processing error:', error);
      message.setReject('Failed to process email');
    }
  },

  // Scheduled tasks (cron jobs)
  async scheduled(event: ScheduledEvent, env: Env, ctx: ExecutionContext): Promise<void> {
    console.log(`⏰ Running scheduled task: ${event.cron}`);
    
    try {
      switch (event.cron) {
        case '0 */6 * * *': // Every 6 hours
          await cleanupOldData(env);
          break;
        
        case '0 0 * * *': // Daily at midnight
          await resetDailyQuotas(env);
          break;
        
        case '*/5 * * * *': // Every 5 minutes
          await checkServiceHealth(env);
          break;
        
        case '0 2 * * *': // Daily at 2 AM
          await backupDatabase(env);
          break;
      }
    } catch (error) {
      console.error('Scheduled task failed:', error);
    }
  },

  // Queue consumer (if using queues)
  async queue(batch: MessageBatch<any>, env: Env, ctx: ExecutionContext): Promise<void> {
    console.log(`📬 Processing queue batch of ${batch.messages.length} messages`);
    
    for (const message of batch.messages) {
      try {
        // Process message
        console.log('Processing message:', message.id);
        
        // Mark as completed
        message.ack();
      } catch (error) {
        console.error('Failed to process message:', message.id, error);
        message.retry();
      }
    }
  },
};

// ============ Utility Functions ============

/**
 * Initialize database if needed
 */
async function initializeDatabase(env: Env): Promise<void> {
  try {
    // Check if users table exists
    const tables = await env.DB.prepare(
      "SELECT name FROM sqlite_master WHERE type='table' AND name='users'"
    ).first();

    if (!tables) {
      console.log('📦 Initializing database...');
      
      // Create tables from schema
      const schema = `
        CREATE TABLE IF NOT EXISTS users (
          id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
          email TEXT UNIQUE NOT NULL,
          username TEXT UNIQUE NOT NULL,
          password_hash TEXT NOT NULL,
          role TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin', 'moderator')),
          email_verified BOOLEAN NOT NULL DEFAULT FALSE,
          quota_text_tokens INTEGER NOT NULL DEFAULT 10000,
          quota_images INTEGER NOT NULL DEFAULT 10,
          quota_videos INTEGER NOT NULL DEFAULT 3,
          quota_audio_minutes INTEGER NOT NULL DEFAULT 30,
          full_name TEXT,
          avatar_url TEXT,
          bio TEXT,
          settings TEXT NOT NULL DEFAULT '{}',
          created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
          updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
          last_login_at DATETIME,
          deleted_at DATETIME
        );

        CREATE TABLE IF NOT EXISTS services (
          id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
          user_id TEXT NOT NULL,
          name TEXT NOT NULL,
          slug TEXT UNIQUE NOT NULL,
          description TEXT,
          status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'active', 'suspended', 'deleted')),
          config TEXT NOT NULL DEFAULT '{}',
          category TEXT,
          tags TEXT,
          icon_url TEXT,
          views INTEGER NOT NULL DEFAULT 0,
          likes INTEGER NOT NULL DEFAULT 0,
          created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
          updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
          approved_at DATETIME,
          deleted_at DATETIME,
          FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS ai_requests (
          id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
          user_id TEXT NOT NULL,
          provider TEXT NOT NULL CHECK (provider IN ('gemini', 'openai', 'huggingface', 'ollama', 'stability')),
          model TEXT NOT NULL,
          type TEXT NOT NULL CHECK (type IN ('text', 'image', 'video', 'audio')),
          input_text TEXT,
          input_tokens INTEGER DEFAULT 0,
          output_text TEXT,
          output_tokens INTEGER DEFAULT 0,
          image_url TEXT,
          video_url TEXT,
          audio_url TEXT,
          cost_usd DECIMAL(10, 4) DEFAULT 0,
          status TEXT NOT NULL DEFAULT 'completed' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
          error_message TEXT,
          metadata TEXT NOT NULL DEFAULT '{}',
          created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
          completed_at DATETIME,
          FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        );

        CREATE TABLE IF NOT EXISTS email_logs (
          id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
          from_email TEXT NOT NULL,
          to_email TEXT NOT NULL,
          subject TEXT,
          message_id TEXT,
          status TEXT NOT NULL DEFAULT 'received' CHECK (status IN ('received', 'forwarded', 'rejected', 'bounced')),
          status_message TEXT,
          headers TEXT,
          body_preview TEXT,
          processed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        );

        CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
        CREATE INDEX IF NOT EXISTS idx_services_user_id ON services(user_id);
        CREATE INDEX IF NOT EXISTS idx_ai_requests_user_id ON ai_requests(user_id);
      `;

      // Execute schema creation
      await env.DB.exec(schema);
      console.log('✅ Database initialized successfully');
    }
  } catch (error) {
    console.error('❌ Database initialization failed:', error);
  }
}

/**
 * Cleanup old data
 */
async function cleanupOldData(env: Env): Promise<void> {
  try {
    // Delete AI requests older than 30 days
    await env.DB.prepare(
      'DELETE FROM ai_requests WHERE created_at < datetime(CURRENT_TIMESTAMP, ?)'
    ).bind('-30 days').run();

    // Delete email logs older than 90 days
    await env.DB.prepare(
      'DELETE FROM email_logs WHERE processed_at < datetime(CURRENT_TIMESTAMP, ?)'
    ).bind('-90 days').run();

    // Soft delete users who haven't logged in for 1 year
    await env.DB.prepare(
      `UPDATE users 
       SET deleted_at = CURRENT_TIMESTAMP 
       WHERE last_login_at < datetime(CURRENT_TIMESTAMP, ?) 
       AND deleted_at IS NULL`
    ).bind('-365 days').run();

    console.log('✅ Old data cleaned up');
  } catch (error) {
    console.error('❌ Cleanup failed:', error);
  }
}

/**
 * Reset daily quotas
 */
async function resetDailyQuotas(env: Env): Promise<void> {
  try {
    // Reset free tier quotas daily
    await env.DB.exec(`
      UPDATE users 
      SET quota_text_tokens = 10000,
          quota_images = 10,
          quota_videos = 3,
          quota_audio_minutes = 30
      WHERE role = 'user' 
      AND deleted_at IS NULL
    `);

    console.log('✅ Daily quotas reset');
  } catch (error) {
    console.error('❌ Quota reset failed:', error);
  }
}

/**
 * Check service health
 */
async function checkServiceHealth(env: Env): Promise<void> {
  try {
    // Check database connection
    await env.DB.prepare('SELECT 1').run();
    
    // Log health status
    await env.KV.put('last_health_check', Date.now().toString(), {
      expirationTtl: 3600,
    });

    console.log('✅ All services healthy');
  } catch (error) {
    console.error('❌ Service health check failed:', error);
  }
}

/**
 * Backup database (simplified - in production use proper backup)
 */
async function backupDatabase(env: Env): Promise<void> {
  try {
    // Create backup in R2
    const backupData = {
      timestamp: new Date().toISOString(),
      tables: ['users', 'services', 'ai_requests', 'email_logs'],
    };

    // For now, just log backup
    console.log('💾 Database backup triggered:', backupData.timestamp);
    
    // In production, you would:
    // 1. Export database
    // 2. Upload to R2
    // 3. Notify admin
    
  } catch (error) {
    console.error('❌ Backup failed:', error);
  }
}

/**
 * Generate API documentation
 */
function generateAPIDocs(): any {
  return {
    api: {
      version: '1.0.0',
      endpoints: {
        health: {
          GET: '/health',
          description: 'Health check endpoint',
        },
        auth: {
          POST: {
            register: '/api/v1/auth/register',
            login: '/api/v1/auth/login',
          },
          GET: {
            me: '/api/v1/auth/me',
          },
        },
        users: {
          GET: {
            list: '/api/v1/users',
            single: '/api/v1/users/:id',
          },
        },
        services: {
          GET: '/api/v1/services',
          POST: '/api/v1/services',
          GET_single: '/api/v1/services/:id',
          PUT: '/api/v1/services/:id',
          DELETE: '/api/v1/services/:id',
        },
        ai: {
          POST: '/api/v1/ai/generate',
          GET: {
            quota: '/api/v1/ai/quota',
            requests: '/api/v1/ai/requests',
          },
        },
        email: {
          GET: {
            logs: '/api/v1/email/logs',
            config: '/api/v1/email/config',
          },
        },
      },
      webhooks: {
        email: {
          POST: '/webhooks/email',
        },
        stripe: {
          POST: '/webhooks/stripe',
        },
      },
    },
  };
}

// ============ Type Definitions ============

// Extend Request type
declare global {
  interface IRequest {
    user?: User;
    params?: Record<string, string>;
    query?: Record<string, string>;
  }
}

// Response helper
const createResponse = (
  data: any,
  status = 200,
  headers: Record<string, string> = {}
): Response => {
  const defaultHeaders = {
    'Content-Type': 'application/json',
    'Access-Control-Allow-Origin': '*',
  };

  return new Response(JSON.stringify(data, null, 2), {
    status,
    headers: { ...defaultHeaders, ...headers },
  });
};

// Cache helper
const cacheResponse = async (
  key: string,
  response: Response,
  env: Env,
  ttl = 3600
): Promise<void> => {
  try {
    await env.KV.put(key, await response.clone().text(), {
      expirationTtl: ttl,
      metadata: {
        cachedAt: Date.now(),
        contentType: response.headers.get('content-type'),
      },
    });
  } catch (error) {
    console.error('Cache error:', error);
  }
};

// Rate limiting helper
const rateLimit = async (
  key: string,
  env: Env,
  limit = 100,
  window = 900
): Promise<{ allowed: boolean; remaining: number }> => {
  const current = Date.now();
  const windowStart = current - window * 1000;

  try {
    // Get current count
    const countKey = `rate_limit:${key}`;
    const count = parseInt((await env.KV.get(countKey)) || '0');

    // Check if exceeded
    if (count >= limit) {
      return { allowed: false, remaining: 0 };
    }

    // Increment count
    await env.KV.put(countKey, (count + 1).toString(), {
      expirationTtl: window,
    });

    return { allowed: true, remaining: limit - count - 1 };
  } catch (error) {
    console.error('Rate limit error:', error);
    return { allowed: true, remaining: limit };
  }
};