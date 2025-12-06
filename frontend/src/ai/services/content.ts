import { aiService, type AIRequest } from './api';

export const contentService = {
  async generateBlogPost(topic: string, options?: any) {
    const request: AIRequest = {
      prompt: `Write a blog post about: ${topic}`,
      options
    };
    return aiService.generateContent(request);
  },

  async generateSocialMediaPost(platform: string, topic: string, options?: any) {
    const request: AIRequest = {
      prompt: `Write a ${platform} post about: ${topic}`,
      options
    };
    return aiService.generateContent(request);
  },

  async translateText(text: string, targetLanguage: string, sourceLanguage?: string) {
    const request: AIRequest = {
      prompt: `Translate: ${text}`,
      options: { targetLanguage, sourceLanguage }
    };
    return aiService.generateContent(request);
  },

  async generateProductDescription(product: string, options?: any) {
    const request: AIRequest = {
      prompt: `Write product description for: ${product}`,
      options
    };
    return aiService.generateContent(request);
  }
};
