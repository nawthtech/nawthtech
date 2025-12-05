import & useState, useCallback, useRef } from 'react';
import { aiService, type AIRequest, type AIResponse } from '../services/api'; 
import { contentService } from '../ services/content';
import { analysisService } from'../ services/analysis';
import { mediaService } from '../ services/media';

interface UseAIOptions {
  onSuccess?: (data: any) => void;
  onError?: (error: Error) => void;
  showNotifications?: boolean;
}

export const useAI = (options: UseAIOptions = {}) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<any>(null);
  const [progress, setProgress] = useState(0);
  
  const abortControllerRef = useRef<AbortController | null>(null);
  
  // توليد محتوى نصي عام
  const generateContent = useCallback(async (request: AIRequest) => {
    setLoading(true);
    setError(null);
    setProgress(0);
    
    abortControllerRef.current = new AbortController();
    
    try {
      // محاكاة التقدم
      const progressInterval = setInterval(() => {
        setProgress(prev => Math.min(prev + 10, 90));
      }, 500);
      
      const response = await aiService.generateContent(request);
      
      clearInterval(progressInterval);
      setProgress(100);
      setResult(response);
      
      if (options.onSuccess) {
        options.onSuccess(response);
      }
      
      return response;
    } catch (err: any) {
      setError(err.message);
      if (options.onError) {
        options.onError(err);
      }
      throw err;
    } finally {
      setLoading(false);
      setTimeout(() => setProgress(0), 1000);
    }
  }, [options]);
  
  // توليد مقال
  const generateBlogPost = useCallback(async (topic: string, language: string = 'ar') => {
    return await contentService.generateBlogPost(topic, { language });
  }, []);
  
  // توليد منشور وسائط اجتماعية
  const generateSocialMediaPost = useCallback(async (
    platform: 'twitter' | 'linkedin' | 'instagram' | 'facebook',
    topic: string,
    language: string = 'ar'
  ) => {
    return await contentService.generateSocialMediaPost(platform, topic, { language });
  }, []);
  
  // تحليل اتجاهات السوق
  const analyzeMarketTrends = useCallback(async (industry: string, timeframe: string) => {
    return await analysisService.analyzeMarketTrends(industry, timeframe);
  }, []);
  
  // توليد صورة
  const generateImage = useCallback(async (prompt: string, style: string = 'realistic') => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await mediaService.generateSocialMediaImage('instagram', prompt, style);
      setResult(response);
      return response;
    } catch (err: any) {
      setError(err.message);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);
  
  // إلغاء العملية الجارية
  const cancel = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
      setLoading(false);
      setError('Operation cancelled');
    }
  }, []);
  
  // إعادة التعيين
  const reset = useCallback(() => {
    setLoading(false);
    setError(null);
    setResult(null);
    setProgress(0);
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
  }, []);
  
  // الحصول على النماذج المتاحة
  const getAvailableModels = useCallback(async () => {
    try {
      return await aiService.getAvailableModels();
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  }, []);
  
  // الحصول على الاستخدام
  const getUsage = useCallback(async () => {
    try {
      return await aiService.getUsage();
    } catch (err: any) {
      setError(err.message);
      throw err;
    }
  }, []);
  
  return {
    // State
    loading,
    error,
    result,
    progress,
    
    // Actions
    generateContent,
    generateBlogPost,
    generateSocialMediaPost,
    analyzeMarketTrends,
    generateImage,
    getAvailableModels,
    getUsage,
    
    // Control
    cancel,
    reset,
    
    // Services
    contentService,
    analysisService,
    mediaService,
  };
};