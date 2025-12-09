-- Migration: 0001_initial
-- Description: Initial database schema
-- Created: 2024-12-09

-- Enable foreign keys
PRAGMA foreign_keys = ON;

-- Users table
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

-- Services table
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

-- AI Requests table
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

-- Email logs table
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

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_created ON users(created_at);

CREATE INDEX IF NOT EXISTS idx_services_user ON services(user_id);
CREATE INDEX IF NOT EXISTS idx_services_slug ON services(slug);
CREATE INDEX IF NOT EXISTS idx_services_status ON services(status);
CREATE INDEX IF NOT EXISTS idx_services_created ON services(created_at);

CREATE INDEX IF NOT EXISTS idx_ai_requests_user ON ai_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_ai_requests_provider ON ai_requests(provider);
CREATE INDEX IF NOT EXISTS idx_ai_requests_type ON ai_requests(type);
CREATE INDEX IF NOT EXISTS idx_ai_requests_created ON ai_requests(created_at);

CREATE INDEX IF NOT EXISTS idx_email_logs_processed ON email_logs(processed_at);

-- Triggers for updated_at
CREATE TRIGGER IF NOT EXISTS update_users_timestamp 
AFTER UPDATE ON users
BEGIN
  UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_services_timestamp 
AFTER UPDATE ON services
BEGIN
  UPDATE services SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Insert default admin user
INSERT OR IGNORE INTO users (id, email, username, password_hash, role, email_verified, full_name)
VALUES (
  '00000000-0000-0000-0000-000000000001',
  'admin@nawthtech.com',
  'admin',
  -- Hash of "Admin@123"
  '8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918',
  'admin',
  TRUE,
  'System Administrator'
);

-- Insert default test user
INSERT OR IGNORE INTO users (id, email, username, password_hash, email_verified, full_name)
VALUES (
  '00000000-0000-0000-0000-000000000002',
  'test@nawthtech.com',
  'testuser',
  -- Hash of "Test@123"
  '8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918',
  TRUE,
  'Test User'
);

-- Migration complete
SELECT '✅ Migration 0001_initial completed successfully!' as message;