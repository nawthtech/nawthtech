import { useState, useCallback, useRef } from 'react';
import { aiService, type AIRequest, type AIResponse } from '../services/api'; 

type Language = 'ar' | 'en' | 'fr' | 'es';

interface UseAIOptions {
  onSuccess?: (data: any) => void;
  onError?: (error: Error) => void;
  showNotifications?: boolean;
}

export const useAI = (options: UseAIOptions = {}) => {
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<any>(null);
  const [progress, setProgress] = useState<number>(0);
  
  const abortControllerRef = useRef<AbortController | null>(null);
  
  const generateContent = useCallback(async (request: AIRequest): Promise<AIResponse> => {
    setLoading(true);
    setError(null);
    setProgress(0);
    
    abortControllerRef.current = new AbortController();
    
    try {
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
      setError(err.message || 'An error occurred');
      if (options.onError) {
        options.onError(err);
      }
      throw err;
    } finally {
      setLoading(false);
      setTimeout(() => setProgress(0), 1000);
    }
  }, [options]);
  
  const generateBlogPost = useCallback(async (topic: string, language: string = 'ar'): Promise<any> => {
    setLoading(true);
    setError(null);
    
    try {
      const request: AIRequest = {
        prompt: `Write a blog post about: ${topic}`,
        options: { language }
      };
      
      const response = await generateContent(request);
      return response;
    } catch (err: any) {
      setError(err.message || 'Failed to generate blog post');
      if (options.onError) {
        options.onError(err);
      }
      throw err;
    } finally {
      setLoading(false);
    }
  }, [generateContent, options.onError]);
  
  const generateSocialMediaPost = useCallback(async (
    platform: 'twitter' | 'linkedin' | 'instagram' | 'facebook',
    topic: string,
    language: string = 'ar'
  ): Promise<any> => {
    setLoading(true);
    setError(null);
    
    try {
      const request: AIRequest = {
        prompt: `Write a ${platform} post about: ${topic}`,
        options: { language }
      };
      
      const response = await generateContent(request);
      return response;
    } catch (err: any) {
      setError(err.message || `Failed to generate ${platform} post`);
      if (options.onError) {
        options.onError(err);
      }
      throw err;
    } finally {
      setLoading(false);
    }
  }, [generateContent, options.onError]);
  
  const analyzeMarketTrends = useCallback(async (industry: string, timeframe: string): Promise<any> => {
    setLoading(true);
    setError(null);
    
    try {
      const request: AIRequest = {
        prompt: `Analyze market trends for ${industry} industry for ${timeframe}`,
      };
      
      const response = await generateContent(request);
      return response;
    } catch (err: any) {
      setError(err.message || 'Failed to analyze market trends');
      if (options.onError) {
        options.onError(err);
      }
      throw err;
    } finally {
      setLoading(false);
    }
  }, [generateContent, options.onError]);
  
  const generateImage = useCallback(async (prompt: string, style: string = 'realistic'): Promise<any> => {
    setLoading(true);
    setError(null);
    
    try {
      const request: AIRequest = {
        prompt: `Generate an image: ${prompt} in ${style} style`,
        options: { style }
      };
      
      const response = await generateContent(request);
      return response;
    } catch (err: any) {
      setError(err.message || 'Failed to generate image');
      if (options.onError) {
        options.onError(err);
      }
      throw err;
    } finally {
      setLoading(false);
    }
  }, [generateContent, options.onError]);
  
  const cancel = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
      setLoading(false);
      setError('Operation cancelled');
    }
  }, []);
  
  const reset = useCallback(() => {
    setLoading(false);
    setError(null);
    setResult(null);
    setProgress(0);
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
    }
  }, []);
  
  const getAvailableModels = useCallback(async (): Promise<any> => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await aiService.getAvailableModels();
      return response;
    } catch (err: any) {
      setError(err.message || 'Failed to get available models');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);
  
  const getUsage = useCallback(async (): Promise<any> => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await aiService.getUsage();
      return response;
    } catch (err: any) {
      setError(err.message || 'Failed to get usage data');
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);
  
  const translateText = useCallback(async (
    text: string,
    targetLanguage: Language,
    sourceLanguage?: Language
  ): Promise<any> => {
    setLoading(true);
    setError(null);
    
    try {
      const request: AIRequest = {
        prompt: `Translate: ${text}`,
        options: { 
          targetLanguage, 
          sourceLanguage: sourceLanguage || 'auto' 
        }
      };
      
      const response = await generateContent(request);
      return response;
    } catch (err: any) {
      setError(err.message || 'Failed to translate text');
      if (options.onError) {
        options.onError(err);
      }
      throw err;
    } finally {
      setLoading(false);
    }
  }, [generateContent, options.onError]);
  
  return {
    loading,
    error,
    result,
    progress,
    generateContent,
    generateBlogPost,
    generateSocialMediaPost,
    analyzeMarketTrends,
    generateImage,
    translateText,
    getAvailableModels,
    getUsage,
    cancel,
    reset,
  };
};
