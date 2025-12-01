/**
 * Environment Variables Type Definitions for Vite
 */

/// <reference types="vite/client" />

interface ImportMetaEnv {
  // Application
  readonly VITE_APP_NAME: string
  readonly VITE_APP_VERSION: string
  readonly VITE_APP_BASE_URL: string
  readonly VITE_APP_ENVIRONMENT: 'development' | 'test' | 'production'
  
  // API
  readonly VITE_API_URL: string
  readonly VITE_API_TIMEOUT: string
  readonly VITE_API_RETRY_ATTEMPTS: string
  
  // Storage
  readonly VITE_STORAGE_URL: string
  readonly VITE_UPLOAD_ENDPOINT: string
  
  // Analytics
  readonly VITE_GA_TRACKING_ID: string
  readonly VITE_FB_PIXEL_ID: string
  readonly VITE_HOTJAR_ID: string
  readonly VITE_MIXPANEL_TOKEN: string
  
  // Social Media
  readonly VITE_INSTAGRAM_APP_ID: string
  readonly VITE_INSTAGRAM_REDIRECT_URI: string
  readonly VITE_TWITTER_API_KEY: string
  readonly VITE_TWITTER_API_SECRET: string
  readonly VITE_FACEBOOK_APP_ID: string
  readonly VITE_LINKEDIN_CLIENT_ID: string
  
  // AI Services
  readonly VITE_OPENAI_API_KEY: string
  readonly VITE_STABILITYAI_API_KEY: string
  readonly VITE_HUGGINGFACE_API_KEY: string
  
  // Payments
  readonly VITE_STRIPE_PUBLISHABLE_KEY: string
  readonly VITE_STRIPE_SECRET_KEY: string
  
  // Notifications
  readonly VITE_VAPID_PUBLIC_KEY: string
  
  // Development
  readonly VITE_USE_MOCK_DATA: string
  readonly VITE_ENABLE_REDUX_DEVTOOLS: string
  readonly VITE_LOG_LEVEL: 'debug' | 'info' | 'warn' | 'error'
  
  // SEO
  readonly VITE_DEFAULT_TITLE: string
  readonly VITE_DEFAULT_DESCRIPTION: string
  readonly VITE_DEFAULT_KEYWORDS: string
  readonly VITE_CANONICAL_URL: string
  
  // Security
  readonly VITE_CSRF_HEADER_NAME: string
  readonly VITE_ALLOWED_ORIGINS: string
  
  // Performance
  readonly VITE_CACHE_DURATION: string
  readonly VITE_LAZY_LOADING: string
  readonly VITE_IMAGE_OPTIMIZATION: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}