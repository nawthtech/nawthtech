import axios, { type AxiosInstance, type AxiosRequestConfig } from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface AIRequest {
  prompt: string;
  model?: string;
  language?: string;
  tone?: string;
  length?: string;
  options?: Record<string, any>;
}

export interface AIResponse {
  success: boolean;
  data: {
    content: string;
    tokens_used?: number;
    model_used: string;
  };
  error?: string;
}

export interface MediaGenerationRequest {
  type: 'image' | 'video' | 'audio';
  prompt: string;
  style?: string;
  dimensions?: { width: number; height: number };
  duration?: number;
}

export interface MediaGenerationResponse {
  success: boolean;
  data: {
    url: string;
    media_type: string;
    size: number;
    duration?: number;
  };
}

class AIService {
  private axiosInstance: AxiosInstance;
  
  constructor() {
    this.axiosInstance = axios.create({
      baseURL: API_BASE_URL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    // إضافة interceptor للـ JWT token
    this.axiosInstance.interceptors.request.use((config: AxiosRequestConfig)) => {
      const token = localStorage.getItem('access_token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });
  }
  
  // توليد محتوى نصي
  async generateContent(request: AIRequest): Promise<AIResponse> {
    try {
      const response = await this.axiosInstance.post('/api/ai/generate', request);
      return response.data;
    } catch (error: any) {
      throw this.handleError(error);
    }
  }
  
  // تحليل الاتجاهات
  async analyzeTrends(industry: string, timeframe: string): Promise<AIResponse> {
    try {
      const response = await this.axiosInstance.post('/api/ai/analyze/trends', {
        industry,
        timeframe,
      });
      return response.data;
    } catch (error: any) {
      throw this.handleError(error);
    }
  }
  
  // توليد استراتيجية
  async generateStrategy(businessType: string, goals: string): Promise<AIResponse> {
    try {
      const response = await this.axiosInstance.post('/api/ai/strategy', {
        business_type: businessType,
        goals,
      });
      return response.data;
    } catch (error: any) {
      throw this.handleError(error);
    }
  }
  
  // توليد صور
  async generateImage(prompt: string, style: string = 'realistic'): Promise<MediaGenerationResponse> {
    try {
      const response = await this.axiosInstance.post('/api/ai/generate/image', {
        prompt,
        style,
        dimensions: { width: 1024, height: 1024 },
      });
      return response.data;
    } catch (error: any) {
      throw this.handleError(error);
    }
  }
  
  // توليد فيديو
  async generateVideo(prompt: string, duration: number = 3): Promise<MediaGenerationResponse> {
    try {
      const response = await this.axiosInstance.post('/api/ai/generate/video', {
        prompt,
        duration,
        style: 'animated',
      });
      return response.data;
    } catch (error: any) {
      throw this.handleError(error);
    }
  }
  
  // ترجمة نص
  async translateText(text: string, targetLang: string): Promise<AIResponse> {
    try {
      const response = await this.axiosInstance.post('/api/ai/translate', {
        text,
        target_language: targetLang,
      });
      return response.data;
    } catch (error: any) {
      throw this.handleError(error);
    }
  }
  
  // الحصول على النماذج المتاحة
  async getAvailableModels(): Promise<any> {
    try {
      const response = await this.axiosInstance.get('/api/ai/models');
      return response.data;
    } catch (error: any) {
      throw this.handleError(error);
    }
  }
  
  // الحصول على الاستخدام والحصص
  async getUsage(): Promise<any> {
    try {
      const response = await this.axiosInstance.get('/api/ai/usage');
      return response.data;
    } catch (error: any) {
      throw this.handleError(error);
    }
  }
  
  private handleError(error: unknown): Error {
    if (axios.isAxiosError(error)) {
      return new Error(error.response?.data?.error || 'AI service error');
    }
    return new Error('Unknown error occurred');
  }
}

export const aiService = new AIService();