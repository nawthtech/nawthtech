import { aiService, type AIRequest } from './api';

export interface ContentGenerationOptions {
  contentType: 'blog_post' | 'social_media' | 'email' | 'ad_copy' | 'product_description';
  tone: 'professional' | 'casual' | 'persuasive' | 'informative';
  language: 'ar' | 'en' | 'fr' | 'es';
  length: 'short' | 'medium' | 'long';
  keywords?: string[];
}

export class ContentService {
  // توليد مقال
  async generateBlogPost(topic: string, options: Partial<ContentGenerationOptions> = {}) {
    const prompt = this.buildBlogPostPrompt(topic, options);
    
    const request: AIRequest = {
      prompt,
      language: options.language || 'ar',
      tone: options.tone || 'professional',
      options: {
        content_type: 'blog_post',
        length: options.length || 'medium',
      },
    };
    
    return await aiService.generateContent(request);
  }
  
  // توليد منشور وسائط اجتماعية
  async generateSocialMediaPost(
    platform: 'twitter' | 'linkedin' | 'instagram' | 'facebook',
    topic: string,
    options: Partial<ContentGenerationOptions> = {}
  ) {
    const prompt = this.buildSocialMediaPrompt(platform, topic, options);
    
    const request: AIRequest = {
      prompt,
      language: options.language || 'ar',
      tone: options.tone || 'casual',
      options: {
        content_type: 'social_media',
        platform,
      },
    };
    
    return await aiService.generateContent(request);
  }
  
  // توليد نص إعلاني
  async generateAdCopy(product: string, targetAudience: string, options: Partial<ContentGenerationOptions> = {}) {
    const prompt = this.buildAdCopyPrompt(product, targetAudience, options);
    
    const request: AIRequest = {
      prompt,
      language: options.language || 'ar',
      tone: options.tone || 'persuasive',
      options: {
        content_type: 'ad_copy',
      },
    };
    
    return await aiService.generateContent(request);
  }
  
  // توليد وصف منتج
  async generateProductDescription(product: string, features: string[], options: Partial<ContentGenerationOptions> = {}) {
    const prompt = this.buildProductDescriptionPrompt(product, features, options);
    
    const request: AIRequest = {
      prompt,
      language: options.language || 'ar',
      tone: options.tone || 'professional',
      options: {
        content_type: 'product_description',
      },
    };
    
    return await aiService.generateContent(request);
  }
  
  // بناء prompts
  private buildBlogPostPrompt(topic: string, options: Partial<ContentGenerationOptions>): string {
    return `اكتب مقالاً عن "${topic}" 
    
    المتطلبات:
    - اللغة: ${options.language === 'ar' ? 'عربية فصحى' : options.language}
    - النبرة: ${options.tone}
    - الطول: ${options.length === 'short' ? 'قصير (300-500 كلمة)' : options.length === 'long' ? 'طويل (1000-1500 كلمة)' : 'متوسط (600-800 كلمة)'}
    - موجه لرواد الأعمال والشركات الناشئة
    - يتضمن نصائح عملية قابلة للتطبيق
    - يحتوي على أمثلة واقعية
    
    التنسيق: Markdown مع عناوين رئيسية وفرعية`;
  }
  
  private buildSocialMediaPrompt(
    platform: string,
    topic: string,
    options: Partial<ContentGenerationOptions>
  ): string {
    const platformInfo: Record<string, string> = {
      twitter: 'تغريدة على تويتر (مختصرة، تصل إلى 280 حرفاً، تتضمن هاشتاقات)',
      linkedin: 'منشور على لينكدإن (مهني، يركز على الصناعة والأعمال)',
      instagram: 'منشور على إنستغرام (جذاب بصرياً، يتضمن وصفاً وإيموجيات)',
      facebook: 'منشور على فيسبوك (اجتماعي، يحفز التفاعل)',
    };
    
    return `اكتب منشوراً لـ${platformInfo[platform]} عن "${topic}"
    
    المتطلبات:
    - اللغة: ${options.language === 'ar' ? 'عربية' : options.language}
    - النبرة: ${options.tone}
    - يتضمن دعوة للتفاعل
    - مناسب للجمهور العربي
    - ${platform === 'twitter' ? 'يتضمن 2-3 هاشتاقات ذات صلة' : ''}
    - ${platform === 'instagram' ? 'يتضمن إيموجيات مناسبة' : ''}
    
    قدم النص جاهزاً للنشر.`;
  }
  
  private buildAdCopyPrompt(product: string, targetAudience: string, options: Partial<ContentGenerationOptions>): string {
    return `اكتب نصاً إعلانياً مقنعاً عن "${product}"
    
    الجمهور المستهدف: ${targetAudience}
    
    المتطلبات:
    - اللغة: ${options.language === 'ar' ? 'عربية مؤثرة' : options.language}
    - النبرة: ${options.tone}
    - يركز على فوائد المنتج
    - يتضمن دعوة واضحة للعمل
    - يستخدم لغة إقناعية
    - يتناول احتياجات الجمهور المستهدف
    
    الطول: 100-200 كلمة`;
  }
  
  private buildProductDescriptionPrompt(
    product: string,
    features: string[],
    options: Partial<ContentGenerationOptions>
  ): string {
    return `اكتب وصفاً تفصيلياً للمنتج: "${product}"
    
    الميزات الرئيسية:
    ${features.map(f => `- ${f}`).join('\n')}
    
    المتطلبات:
    - اللغة: ${options.language === 'ar' ? 'عربية واضحة' : options.language}
    - النبرة: ${options.tone}
    - يركز على فوائد المنتج للعميل
    - يستخدم لغة جذابة
    - يتضمن معلومات تقنية بشكل مبسط
    - مناسب لموقع إلكتروني أو متجر إلكتروني
    
    التنسيق: فقرات واضحة مع عناوين فرعية`;
  }
}

export const contentService = new ContentService();