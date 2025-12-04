import { aiService } from './api';

export interface AnalysisOptions {
  depth: 'basic' | 'detailed' | 'comprehensive';
  format: 'report' | 'summary' | 'bullet_points';
  includeRecommendations: boolean;
}

export class AnalysisService {
  // تحليل اتجاهات السوق
  async analyzeMarketTrends(industry: string, timeframe: string, options: Partial<AnalysisOptions> = {}) {
    const prompt = `قم بتحليل اتجاهات السوق لصناعة ${industry} خلال ${timeframe}
    
    عمق التحليل: ${options.depth || 'detailed'}
    
    أدرج في التحليل:
    1. حجم السوق الحالي ومعدل النمو
    2. المحركات الرئيسية والتحديات
    3. التقنيات الناشئة
    4. المشهد التنافسي
    5. التوقعات المستقبلية
    ${options.includeRecommendations ? '6. توصيات استراتيجية للشركات' : ''}
    
    التنسيق: ${options.format === 'report' ? 'تقرير منظم' : options.format === 'bullet_points' ? 'نقاط رئيسية' : 'ملخص تنفيذي'}`;
    
    return await aiService.generateContent({ prompt });
  }
  
  // تحليل SWOT
  async performSWOTAnalysis(
    companyName: string,
    industry: string,
    options: Partial<AnalysisOptions> = {}
  ) {
    const prompt = `قم بإجراء تحليل SWOT لشركة ${companyName} في صناعة ${industry}
    
    قم بتضمين:
    
    نقاط القوة (الداخلية):
    - اذكر 5-7 نقاط قوة
    
    نقاط الضعف (الداخلية):
    - اذكر 5-7 نقاط ضعف
    
    الفرص (الخارجية):
    - اذكر 5-7 فرص
    
    التهديدات (الخارجية):
    - اذكر 5-7 تهديدات
    
    ${options.includeRecommendations ? 'التوصيات الاستراتيجية:' : ''}
    ${options.includeRecommendations ? '- بناءً على تحليل SWOT' : ''}
    
    كن واقعياً وموضوعياً.`;
    
    return await aiService.generateContent({ prompt });
  }
  
  // تحليل الجمهور المستهدف
  async analyzeTargetAudience(
    demographics: string,
    interests: string,
    behavior: string
  ) {
    const prompt = `قم بتحليل الجمهور المستهدف بناءً على:
    
    الديموغرافيا: ${demographics}
    الاهتمامات: ${interests}
    السلوك: ${behavior}
    
    قدم تحليلاً يتضمن:
    1. تفضيلات المحتوى
    2. استخدام المنصات الرقمية
    3. محفزات الشراء
    4. نقاط الألم
    5. فرص التخصيص
    6. توصيات للتواصل الفعال
    
    ركز على التطبيقات العملية للتسويق الرقمي.`;
    
    return await aiService.generateContent({ prompt });
  }
  
  // تحليل المنافسين
  async analyzeCompetitors(competitors: string[], industry: string) {
    const prompt = `قم بتحليل المنافسين التاليين في صناعة ${industry}:
    ${competitors.map(c => `- ${c}`).join('\n')}
    
    قارن بينهم بناءً على:
    1. نقاط القوة والضعف
    2. استراتيجيات التسويق
    3. حضور رقمي
    4. عروض المنتجات/الخدمات
    5. التميز التنافسي
    
    قدم توصيات لكيفية التفوق عليهم.`;
    
    return await aiService.generateContent({ prompt });
  }
}

export const analysisService = new AnalysisService();