/**
 * Application Constants
 */

// API Constants
export const API_CONSTANTS = {
  STATUS: {
    SUCCESS: 200,
    CREATED: 201,
    BAD_REQUEST: 400,
    UNAUTHORIZED: 401,
    FORBIDDEN: 403,
    NOT_FOUND: 404,
    INTERNAL_ERROR: 500,
  },
  HEADERS: {
    CONTENT_TYPE: 'Content-Type',
    AUTHORIZATION: 'Authorization',
    ACCEPT_LANGUAGE: 'Accept-Language',
    CSRF_TOKEN: 'X-CSRF-Token',
  },
  CONTENT_TYPES: {
    JSON: 'application/json',
    FORM_DATA: 'multipart/form-data',
  },
  TIMEOUT: 30000,
  MAX_RETRIES: 3,
};

// Storage Constants
export const STORAGE_CONSTANTS = {
  KEYS: {
    AUTH_TOKEN: 'auth_token',
    USER_DATA: 'user_data',
    THEME: 'theme',
    LANGUAGE: 'language',
    CART: 'cart',
    SETTINGS: 'app_settings',
  },
  PREFIX: 'nawthtech_',
  VERSION: 'v1',
};

// Feature Flags Constants
export const FEATURE_FLAGS = {
  AI_ASSISTANT: 'aiAssistant',
  SOCIAL_MEDIA: 'socialMediaIntegration',
  ANALYTICS: 'analytics',
  MULTI_LANGUAGE: 'multiLanguage',
  DARK_MODE: 'darkMode',
  PUSH_NOTIFICATIONS: 'pushNotifications',
} as const;

// Theme Constants
export const THEME_CONSTANTS = {
  MODES: {
    LIGHT: 'light',
    DARK: 'dark',
    AUTO: 'auto',
  },
  DIRECTION: {
    RTL: 'rtl',
    LTR: 'ltr',
  },
  COLOR_SCHEMES: ['blue', 'green', 'purple', 'orange'] as const,
};

// Language Constants
export const LANGUAGE_CONSTANTS = {
  DEFAULT: 'ar',
  SUPPORTED: [
    { code: 'ar', name: 'العربية', dir: 'rtl', flag: '🇸🇦' },
    { code: 'en', name: 'English', dir: 'ltr', flag: '🇺🇸' },
  ] as const,
};

// File Upload Constants
export const UPLOAD_CONSTANTS = {
  MAX_SIZE: 100 * 1024 * 1024, // 100MB
  CHUNK_SIZE: 5 * 1024 * 1024, // 5MB
  MAX_FILES: 10,
  IMAGE_TYPES: ['image/jpeg', 'image/jpg', 'image/png', 'image/webp', 'image/gif'],
  DOCUMENT_TYPES: ['application/pdf', 'application/msword', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'],
  VIDEO_TYPES: ['video/mp4', 'video/mpeg', 'video/quicktime'],
};

// Security Constants
export const SECURITY_CONSTANTS = {
  PASSWORD: {
    MIN_LENGTH: 6,
    REQUIRE_UPPERCASE: false,
    REQUIRE_LOWERCASE: false,
    REQUIRE_NUMBERS: true,
    REQUIRE_SPECIAL_CHARS: false,
  },
  SESSION: {
    TIMEOUT: 24 * 60 * 60 * 1000, // 24 hours
    REFRESH_INTERVAL: 60 * 60 * 1000, // 1 hour
  },
};

// Analytics Constants
export const ANALYTICS_CONSTANTS = {
  EVENTS: {
    PAGE_VIEW: 'page_view',
    BUTTON_CLICK: 'button_click',
    FORM_SUBMIT: 'form_submit',
    FILE_UPLOAD: 'file_upload',
    PAYMENT_COMPLETE: 'payment_complete',
  },
  SCREENS: {
    LOGIN: 'login',
    DASHBOARD: 'dashboard',
    SETTINGS: 'settings',
    PROFILE: 'profile',
  },
};

// Social Media Constants
export const SOCIAL_CONSTANTS = {
  PLATFORMS: {
    INSTAGRAM: 'instagram',
    TWITTER: 'twitter',
    FACEBOOK: 'facebook',
    LINKEDIN: 'linkedin',
    WHATSAPP: 'whatsapp',
    TELEGRAM: 'telegram',
  } as const,
};

// Notification Constants
export const NOTIFICATION_CONSTANTS = {
  TYPES: {
    SUCCESS: 'success',
    ERROR: 'error',
    WARNING: 'warning',
    INFO: 'info',
  },
  POSITIONS: {
    TOP_RIGHT: 'top-right',
    TOP_LEFT: 'top-left',
    BOTTOM_RIGHT: 'bottom-right',
    BOTTOM_LEFT: 'bottom-left',
  },
  DURATION: 5000,
};

// Payment Constants
export const PAYMENT_CONSTANTS = {
  CURRENCY: 'SAR',
  TAX_RATE: 0.15,
  PLANS: {
    BASIC: 'basic',
    PRO: 'pro',
    ENTERPRISE: 'enterprise',
  } as const,
};

// Development Constants
export const DEVELOPMENT_CONSTANTS = {
  LOG_LEVELS: {
    DEBUG: 'debug',
    INFO: 'info',
    WARN: 'warn',
    ERROR: 'error',
  } as const,
};