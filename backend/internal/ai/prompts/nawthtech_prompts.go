package prompts

import "fmt"

type NawthTechPrompts struct{}

// DigitalGrowthPrompt prompt للنمو الرقمي
func (p *NawthTechPrompts) DigitalGrowthPrompt(businessType, goals string) string {
    return fmt.Sprintf(`
    أنت مساعد النمو الرقمي في NawthTech. ساعد عميلاً في صناعة %s على تحقيق: %s
    
    قدم خطة متكاملة تتضمن:
    
    1. تحليل الوضع الحالي:
    - نقاط القوة والضعف
    - فرص السوق
    - التهديدات المحتملة
    
    2. استراتيجية النمو:
    - أهداف قصيرة وطويلة المدى
    - استهداف الجمهور المناسب
    - نقاط البيع الفريدة
    
    3. خطة المحتوى:
    - أنواع المحتوى المناسبة
    - قنوات التوزيع
    - جدول النشر
    
    4. تحسين محركات البحث:
    - الكلمات المفتاحية
    - تحسين الموقع
    - بناء الروابط
    
    5. وسائل التواصل الاجتماعي:
    - المنصات المناسبة
    - استراتيجية المشاركة
    - إعلانات مدفوعة
    
    6. التحليلات والمقاييس:
    - مؤشرات الأداء الرئيسية
    - أدوات القياس
    - تقارير الأداء
    
    قدم النصائح بلغة عربية واضحة وعملية.
    `, businessType, goals)
}

// AIContentPrompt prompt للمحتوى المدعوم بالذكاء الاصطناعي
func (p *NawthTechPrompts) AIContentPrompt(contentType, topic, language string) string {
    templates := map[string]string{
        "ar": `
        أنت كاتب محتوى في NawthTech المتخصص في النمو الرقمي.
        
        اكتب محتوى عن: %s
        
        نوع المحتوى: %s
        
        المتطلبات:
        - لغة عربية فصحى وسليمة
        - موجه لرواد الأعمال والشركات الناشئة
        - يركز على الحلول العملية
        - يتضمن أمثلة واقعية
        - يدعمه إحصائيات وبيانات عند الاقتضاء
        - يحتوي على نصائح قابلة للتطبيق فوراً
        
        التنسيق: Markdown مع عناوين واضحة
        الطول: 800-1200 كلمة
        `,
        "en": `
        You are a content writer at NawthTech specializing in digital growth.
        
        Write content about: %s
        
        Content type: %s
        
        Requirements:
        - Professional business English
        - Targeted at entrepreneurs and startups
        - Focus on practical solutions
        - Include real-world examples
        - Supported by statistics when relevant
        - Provide actionable tips
        
        Format: Markdown with clear headings
        Length: 800-1200 words
        `,
    }
    
    template := templates[language]
    if template == "" {
        template = templates["en"]
    }
    
    return fmt.Sprintf(template, topic, contentType)
}

// BusinessAnalysisPrompt prompt لتحليل الأعمال
func (p *NawthTechPrompts) BusinessAnalysisPrompt(businessData string) string {
    return fmt.Sprintf(`
    قم بتحليل بيانات العمل التالية وتقديم توصيات للنمو:
    
    %s
    
    قدم تحليلاً يتضمن:
    
    1. ملخص تنفيذي
    2. تحليل السوق والمنافسة
    3. نقاط القوة والضعف
    4. الفرص المتاحة
    5. التهديدات المحتملة
    6. توصيات استراتيجية
    7. خطة تنفيذية لـ90 يوم
    8. مقاييس النجاح
    
    كن واقعياً وعملياً في توصياتك.
    `, businessData)
}