import { aiService, type AIRequest } from './api';

export const mediaService = {
  async generateSocialMediaImage(platform: string, prompt: string, style?: string, options?: any) {
    const request: AIRequest = {
      prompt: `Generate ${platform} image: ${prompt} ${style ? `in ${style} style` : ''}`,
      options: { ...options, style }
    };
    return aiService.generateContent(request);
  },

  async generateVideo(prompt: string, options?: any) {
    const request: AIRequest = {
      prompt: `Generate video: ${prompt}`,
      options
    };
    return aiService.generateContent(request);
  },

  async generateAudio(prompt: string, options?: any) {
    const request: AIRequest = {
      prompt: `Generate audio: ${prompt}`,
      options
    };
    return aiService.generateContent(request);
  }
};
