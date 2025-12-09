import type { Env } from '../index';

// Database types
export interface User {
  id: string;
  email: string;
  name?: string;
  password_hash: string;
  role: 'user' | 'admin' | 'super_admin';
  status: 'active' | 'inactive' | 'suspended';
  email_verified: boolean;
  created_at: string;
  updated_at: string;
  last_login_at?: string;
  metadata?: Record<string, any>;
}

export interface Service {
  id: string;
  user_id: string;
  name: string;
  description?: string;
  type: 'ai_text' | 'ai_image' | 'ai_video' | 'storage' | 'api';
  status: 'active' | 'inactive' | 'pending';
  quota_used: number;
  quota_limit: number;
  settings: Record<string, any>;
  created_at: string;
  updated_at: string;
}

export interface ApiKey {
  id: string;
  user_id: string;
  name: string;
  key_hash: string;
  prefix: string;
  last_used_at?: string;
  expires_at?: string;
  permissions: string[];
  created_at: string;
}

export class DatabaseService {
  constructor(private env: Env) {}

  // ============ Users ============
  async createUser(user: Omit<User, 'id' | 'created_at' | 'updated_at'>): Promise<User> {
    const id = crypto.randomUUID();
    const now = new Date().toISOString();
    
    const result = await this.env.DB.prepare(`
      INSERT INTO users (id, email, name, password_hash, role, status, email_verified, created_at, updated_at, metadata)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
      RETURNING *
    `)
      .bind(
        id,
        user.email,
        user.name || null,
        user.password_hash,
        user.role || 'user',
        user.status || 'active',
        user.email_verified ? 1 : 0,
        now,
        now,
        user.metadata ? JSON.stringify(user.metadata) : null
      )
      .first<User>();
    
    return result!;
  }

  async getUserById(id: string): Promise<User | null> {
    const result = await this.env.DB.prepare(`
      SELECT * FROM users WHERE id = ?
    `).bind(id).first<User>();
    
    return result;
  }

  async getUserByEmail(email: string): Promise<User | null> {
    const result = await this.env.DB.prepare(`
      SELECT * FROM users WHERE email = ?
    `).bind(email).first<User>();
    
    return result;
  }

  async updateUser(id: string, updates: Partial<User>): Promise<User | null> {
    const setClauses: string[] = [];
    const values: any[] = [];
    
    if (updates.name !== undefined) {
      setClauses.push('name = ?');
      values.push(updates.name);
    }
    if (updates.email !== undefined) {
      setClauses.push('email = ?');
      values.push(updates.email);
    }
    if (updates.status !== undefined) {
      setClauses.push('status = ?');
      values.push(updates.status);
    }
    if (updates.last_login_at !== undefined) {
      setClauses.push('last_login_at = ?');
      values.push(updates.last_login_at);
    }
    if (updates.metadata !== undefined) {
      setClauses.push('metadata = ?');
      values.push(JSON.stringify(updates.metadata));
    }
    
    setClauses.push('updated_at = ?');
    values.push(new Date().toISOString());
    
    values.push(id);
    
    const result = await this.env.DB.prepare(`
      UPDATE users 
      SET ${setClauses.join(', ')} 
      WHERE id = ?
      RETURNING *
    `).bind(...values).first<User>();
    
    return result;
  }

  // ============ Services ============
  async createService(service: Omit<Service, 'id' | 'created_at' | 'updated_at'>): Promise<Service> {
    const id = crypto.randomUUID();
    const now = new Date().toISOString();
    
    const result = await this.env.DB.prepare(`
      INSERT INTO services (id, user_id, name, description, type, status, quota_used, quota_limit, settings, created_at, updated_at)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
      RETURNING *
    `)
      .bind(
        id,
        service.user_id,
        service.name,
        service.description || null,
        service.type,
        service.status || 'active',
        service.quota_used || 0,
        service.quota_limit,
        service.settings ? JSON.stringify(service.settings) : null,
        now,
        now
      )
      .first<Service>();
    
    return result!;
  }

  async getUserServices(userId: string): Promise<Service[]> {
    const result = await this.env.DB.prepare(`
      SELECT * FROM services WHERE user_id = ? ORDER BY created_at DESC
    `).bind(userId).all<Service>();
    
    return result.results || [];
  }

  async updateServiceQuota(serviceId: string, used: number): Promise<Service | null> {
    const result = await this.env.DB.prepare(`
      UPDATE services 
      SET quota_used = quota_used + ?, updated_at = ?
      WHERE id = ?
      RETURNING *
    `)
      .bind(used, new Date().toISOString(), serviceId)
      .first<Service>();
    
    return result;
  }

  // ============ API Keys ============
  async createApiKey(apiKey: Omit<ApiKey, 'id' | 'created_at'>): Promise<ApiKey> {
    const id = crypto.randomUUID();
    const now = new Date().toISOString();
    
    const result = await this.env.DB.prepare(`
      INSERT INTO api_keys (id, user_id, name, key_hash, prefix, permissions, created_at)
      VALUES (?, ?, ?, ?, ?, ?, ?)
      RETURNING *
    `)
      .bind(
        id,
        apiKey.user_id,
        apiKey.name,
        apiKey.key_hash,
        apiKey.prefix,
        JSON.stringify(apiKey.permissions),
        now
      )
      .first<ApiKey>();
    
    return result!;
  }

  async getApiKeyByHash(keyHash: string): Promise<ApiKey | null> {
    const result = await this.env.DB.prepare(`
      SELECT * FROM api_keys WHERE key_hash = ?
    `).bind(keyHash).first<ApiKey>();
    
    return result;
  }

  async updateApiKeyLastUsed(id: string): Promise<void> {
    await this.env.DB.prepare(`
      UPDATE api_keys SET last_used_at = ? WHERE id = ?
    `).bind(new Date().toISOString(), id).run();
  }

  // ============ Analytics ============
  async logRequest(userId: string | null, endpoint: string, method: string, status: number): Promise<void> {
    await this.env.DB.prepare(`
      INSERT INTO request_logs (id, user_id, endpoint, method, status, created_at)
      VALUES (?, ?, ?, ?, ?, ?)
    `)
      .bind(
        crypto.randomUUID(),
        userId,
        endpoint,
        method,
        status,
        new Date().toISOString()
      )
      .run();
  }

  async getUserStats(userId: string): Promise<{
    total_requests: number;
    services_count: number;
    last_activity: string;
  }> {
    const stats = await this.env.DB.prepare(`
      SELECT 
        (SELECT COUNT(*) FROM request_logs WHERE user_id = ?) as total_requests,
        (SELECT COUNT(*) FROM services WHERE user_id = ?) as services_count,
        (SELECT MAX(created_at) FROM request_logs WHERE user_id = ?) as last_activity
    `)
      .bind(userId, userId, userId)
      .first<{
        total_requests: number;
        services_count: number;
        last_activity: string;
      }>();
    
    return stats!;
  }

  // ============ Migrations ============
  async runMigrations(): Promise<void> {
    // Check if migrations table exists
    const { results } = await this.env.DB.prepare(`
      SELECT name FROM sqlite_master WHERE type='table' AND name='migrations'
    `).all();
    
    if (!results || results.length === 0) {
      // Create migrations table
      await this.env.DB.prepare(`
        CREATE TABLE migrations (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          name TEXT NOT NULL UNIQUE,
          applied_at TEXT NOT NULL
        )
      `).run();
    }
    
    // Run migrations
    const migrations = [
      // Initial schema
      `CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        email TEXT NOT NULL UNIQUE,
        name TEXT,
        password_hash TEXT NOT NULL,
        role TEXT NOT NULL DEFAULT 'user',
        status TEXT NOT NULL DEFAULT 'active',
        email_verified BOOLEAN NOT NULL DEFAULT 0,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL,
        last_login_at TEXT,
        metadata TEXT
      )`,
      
      `CREATE TABLE IF NOT EXISTS services (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        name TEXT NOT NULL,
        description TEXT,
        type TEXT NOT NULL,
        status TEXT NOT NULL DEFAULT 'active',
        quota_used INTEGER NOT NULL DEFAULT 0,
        quota_limit INTEGER NOT NULL,
        settings TEXT,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
      )`,
      
      `CREATE TABLE IF NOT EXISTS api_keys (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        name TEXT NOT NULL,
        key_hash TEXT NOT NULL UNIQUE,
        prefix TEXT NOT NULL,
        last_used_at TEXT,
        expires_at TEXT,
        permissions TEXT NOT NULL,
        created_at TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
      )`,
      
      `CREATE TABLE IF NOT EXISTS request_logs (
        id TEXT PRIMARY KEY,
        user_id TEXT,
        endpoint TEXT NOT NULL,
        method TEXT NOT NULL,
        status INTEGER NOT NULL,
        created_at TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL
      )`,
      
      // Indexes
      `CREATE INDEX IF NOT EXISTS idx_users_email ON users (email)`,
      `CREATE INDEX IF NOT EXISTS idx_users_status ON users (status)`,
      `CREATE INDEX IF NOT EXISTS idx_services_user_id ON services (user_id)`,
      `CREATE INDEX IF NOT EXISTS idx_services_type ON services (type)`,
      `CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys (user_id)`,
      `CREATE INDEX IF NOT EXISTS idx_api_keys_prefix ON api_keys (prefix)`,
      `CREATE INDEX IF NOT EXISTS idx_request_logs_user_id ON request_logs (user_id)`,
      `CREATE INDEX IF NOT EXISTS idx_request_logs_created_at ON request_logs (created_at)`,
    ];
    
    for (const migration of migrations) {
      const migrationName = `migration_${Date.now()}_${migration.substring(0, 50)}`;
      
      // Check if migration already applied
      const applied = await this.env.DB.prepare(`
        SELECT 1 FROM migrations WHERE name = ?
      `).bind(migrationName).first();
      
      if (!applied) {
        try {
          await this.env.DB.prepare(migration).run();
          await this.env.DB.prepare(`
            INSERT INTO migrations (name, applied_at) VALUES (?, ?)
          `).bind(migrationName, new Date().toISOString()).run();
          
          console.log(`✅ Applied migration: ${migrationName}`);
        } catch (error) {
          console.error(`❌ Failed to apply migration ${migrationName}:`, error);
        }
      }
    }
  }
}

// Singleton instance
let dbInstance: DatabaseService | null = null;

export function getDatabaseService(env: Env): DatabaseService {
  if (!dbInstance) {
    dbInstance = new DatabaseService(env);
  }
  return dbInstance;
}