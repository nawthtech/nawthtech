import { useState, useCallback } from 'react';
import { contentService, ContentGenerationOptions } from '../services/content';

interface UseContentGenerationOptions {
  autoSave?: boolean;
  saveInterval?: number;
}

export const useContentGeneration = (options: UseContentGenerationOptions = {}) => {
  const [generatedContent, setGeneratedContent] = useState<string>('');
  const [history, setHistory] = useState<Array<{
    id: string;
    content: string;
    type: string;
    timestamp: Date;
  }>>([]);
  const [isGenerating, setIsGenerating] = useState(false);
  
  // توليد محتوى مع حفظ في التاريخ
  const generateAndSave = useCallback(async (
    type: ContentGenerationOptions['contentType'],
    topic: string,
    options: Partial<ContentGenerationOptions> = {}
  ) => {
    setIsGenerating(true);
    
    try {
      let result;
      
      switch (type) {
        case 'blog_post':
          result = await contentService.generateBlogPost(topic, options);
          break;
        case 'social_media':
          // Assuming platform is provided in options
          result = await contentService.generateSocialMediaPost('twitter', topic, options);
          break;
        case 'email':
          // Implement email generation
          break;
        case 'ad_copy':
          result = await contentService.generateAdCopy(topic, 'business owners', options);
          break;
        case 'product_description':
          result = await contentService.generateProductDescription(topic, [], options);
          break;
        default:
          throw new Error('Unsupported content type');
      }
      
      if (result.success) {
        const newContent = result.data.content;
        setGeneratedContent(newContent);
        
        // إضافة إلى التاريخ
        const historyItem = {
          id: Date.now().toString(),
          content: newContent,
          type,
          timestamp: new Date(),
        };
        
        setHistory(prev => [historyItem, ...prev.slice(0, 9)]); // احتفظ بـ 10 آخرين فقط
        
        return newContent;
      }
      
      throw new Error('Generation failed');
    } finally {
      setIsGenerating(false);
    }
  }, []);
  
  // تحرير المحتوى المولد
  const editContent = useCallback((newContent: string) => {
    setGeneratedContent(newContent);
  }, []);
  
  // حفظ المحتوى
  const saveContent = useCallback((title: string, tags: string[] = []) => {
    const contentToSave = {
      id: Date.now().toString(),
      title,
      content: generatedContent,
      tags,
      createdAt: new Date(),
    };
    
    // حفظ في localStorage
    const savedContents = JSON.parse(localStorage.getItem('nawthtech_contents') || '[]');
    savedContents.unshift(contentToSave);
    localStorage.setItem('nawthtech_contents', JSON.stringify(savedContents.slice(0, 50)));
    
    return contentToSave;
  }, [generatedContent]);
  
  // تحميل من التاريخ
  const loadFromHistory = useCallback((id: string) => {
    const item = history.find(item => item.id === id);
    if (item) {
      setGeneratedContent(item.content);
    }
  }, [history]);
  
  // تصدير المحتوى
  const exportContent = useCallback((format: 'txt' | 'md' | 'html' | 'pdf') => {
    let content = generatedContent;
    let mimeType = 'text/plain';
    let extension = 'txt';
    
    switch (format) {
      case 'md':
        mimeType = 'text/markdown';
        extension = 'md';
        break;
      case 'html':
        content = `<html><body>${content}</body></html>`;
        mimeType = 'text/html';
        extension = 'html';
        break;
      case 'pdf':
        // Note: PDF generation would need a library
        break;
    }
    
    const blob = new Blob([content], { type: mimeType });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `nawthtech_content_${Date.now()}.${extension}`;
    a.click();
    URL.revokeObjectURL(url);
  }, [generatedContent]);
  
  // نسخ إلى الحافظة
  const copyToClipboard = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(generatedContent);
      return true;
    } catch (err) {
      console.error('Failed to copy:', err);
      return false;
    }
  }, [generatedContent]);
  
  return {
    // State
    generatedContent,
    history,
    isGenerating,
    
    // Actions
    generateAndSave,
    editContent,
    saveContent,
    loadFromHistory,
    exportContent,
    copyToClipboard,
    
    // Utilities
    clearContent: () => setGeneratedContent(''),
    clearHistory: () => setHistory([]),
  };
};