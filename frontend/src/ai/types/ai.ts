export interface AIModel {
  id: string;
  name: string;
  provider: string;
  capabilities: string[];
  isLocal: boolean;
  isFree: boolean;
  maxTokens?: number;
  languages: string[];
}

export interface AIConfig {
  defaultModel: string;
  maxTokens: number;
  temperature: number;
  language: string;
  autoSave: boolean;
  notifications: boolean;
}

export interface AIUsage {
  text_used: number;
  text_limit: number;
  images_used: number;
  images_limit: number;
  videos_used: number;
  videos_limit: number;
  audio_used: number;
  audio_limit: number;
  reset_date: string;
}

export interface ContentHistoryItem {
  id: string;
  type: string;
  content: string;
  prompt: string;
  model: string;
  timestamp: Date;
  tokens_used: number;
}

export interface MediaItem {
  id: string;
  type: 'image' | 'video' | 'audio';
  url: string;
  prompt: string;
  style: string;
  size: number;
  duration?: number;
  timestamp: Date;
}

export interface GenerationRequest {
  prompt: string;
  model?: string;
  options?: {
    temperature?: number;
    max_tokens?: number;
    top_p?: number;
    frequency_penalty?: number;
    presence_penalty?: number;
  };
}

export interface GenerationResponse {
  success: boolean;
  data: {
    content?: string;
    media_url?: string;
    tokens_used?: number;
    model_used: string;
    duration?: number;
    size?: number;
  };
  error?: string;
  warnings?: string[];
}