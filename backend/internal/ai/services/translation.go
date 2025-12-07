package services

import (
    "context"
    "fmt"
    "strings"
    "time"
    "github.com/nawthtech/nawthtech/backend/internal/ai/types"
)

type TranslationService struct {
    textProvider types.TextProvider
}

func NewTranslationService(provider types.TextProvider) *TranslationService {
    return &TranslationService{
        textProvider: provider,
    }
}

func (s *TranslationService) Translate(ctx context.Context, text string, sourceLang string, targetLang string) (*types.TranslationResponse, error) {
    req := types.TranslationRequest{
        Text:        text,
        FromLang:    sourceLang,
        ToLang:      targetLang,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.TranslateText(req)
}

func (s *TranslationService) BatchTranslate(ctx context.Context, texts []string, sourceLang string, targetLang string) ([]*types.TranslationResponse, error) {
    var responses []*types.TranslationResponse
    
    for _, text := range texts {
        if strings.TrimSpace(text) == "" {
            continue // تخطي النصوص الفارغة
        }
        
        req := types.TranslationRequest{
            Text:        text,
            FromLang:    sourceLang,
            ToLang:      targetLang,
            UserID:      extractUserIDFromContext(ctx),
            UserTier:    extractUserTierFromContext(ctx),
        }
        
        resp, err := s.textProvider.TranslateText(req)
        if err != nil {
            return responses, fmt.Errorf("failed to translate text '%s': %w", text, err)
        }
        
        responses = append(responses, resp)
        
        // تأخير بسيط بين الطلبات
        select {
        case <-ctx.Done():
            return responses, ctx.Err()
        case <-time.After(50 * time.Millisecond):
            // استمر
        }
    }
    
    return responses, nil
}

func (s *TranslationService) DetectLanguage(ctx context.Context, text string) (string, float64, error) {
    // يمكن استخدام نموذج للكشف عن اللغة أو خدمة خارجية
    // هنا نستخدم نهج بسيط
    prompt := fmt.Sprintf("Detect the language of this text and return only the language code: %s", text)
    
    textReq := types.TextRequest{
        Prompt:   prompt,
        MaxTokens: 10,
        UserID:   extractUserIDFromContext(ctx),
        UserTier: extractUserTierFromContext(ctx),
    }
    
    // نستخدم GenerateText لأن TranslateText يتطلب لغة المصدر
    resp, err := s.textProvider.GenerateText(textReq)
    if err != nil {
        return "", 0.0, fmt.Errorf("failed to detect language: %w", err)
    }
    
    detectedLang := strings.TrimSpace(strings.ToLower(resp.Text))
    
    // محاولة تحديد الثقة (هذه قيمة افتراضية)
    confidence := 0.85
    
    return detectedLang, confidence, nil
}

func (s *TranslationService) TranslateDocument(ctx context.Context, document string, sourceLang string, targetLang string, format string) (*types.TranslationResponse, error) {
    // معالجة خاصة للمستندات حسب التنسيق
    var translatedText string
    
    switch format {
    case "markdown", "md":
        translatedText = s.translateMarkdown(ctx, document, sourceLang, targetLang)
    case "html":
        translatedText = s.translateHTML(ctx, document, sourceLang, targetLang)
    case "json":
        translatedText = s.translateJSON(ctx, document, sourceLang, targetLang)
    default:
        // ترجمة عادية
        req := types.TranslationRequest{
            Text:        document,
            FromLang:    sourceLang,
            ToLang:      targetLang,
            UserID:      extractUserIDFromContext(ctx),
            UserTier:    extractUserTierFromContext(ctx),
        }
        
        resp, err := s.textProvider.TranslateText(req)
        if err != nil {
            return nil, err
        }
        translatedText = resp.TranslatedText
    }
    
    return &types.TranslationResponse{
        TranslatedText: translatedText,
        Cost:           0.0, // سيتم تحديثه من قبل المزود
        Model:          "document-translation",
        CreatedAt:      time.Now(),
    }, nil
}

func (s *TranslationService) translateMarkdown(ctx context.Context, markdown string, sourceLang string, targetLang string) string {
    // تقسيم الماركداون إلى أجزاء (عناوين، فقرات، قوائم)
    lines := strings.Split(markdown, "\n")
    var translatedLines []string
    
    for _, line := range lines {
        if strings.TrimSpace(line) == "" {
            translatedLines = append(translatedLines, "")
            continue
        }
        
        // الحفاظ على تنسيق الماركداون (#, *, `, etc.)
        if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "`") {
            // ترجمة المحتوى فقط، مع الحفاظ على التنسيق
            content := extractContentFromMarkdown(line)
            if content != "" {
                req := types.TranslationRequest{
                    Text:        content,
                    FromLang:    sourceLang,
                    ToLang:      targetLang,
                    UserID:      extractUserIDFromContext(ctx),
                    UserTier:    extractUserTierFromContext(ctx),
                }
                
                if resp, err := s.textProvider.TranslateText(req); err == nil {
                    translatedLine := applyTranslationToMarkdown(line, content, resp.TranslatedText)
                    translatedLines = append(translatedLines, translatedLine)
                    continue
                }
            }
        }
        
        // ترجمة السطر كاملاً
        req := types.TranslationRequest{
            Text:        line,
            FromLang:    sourceLang,
            ToLang:      targetLang,
            UserID:      extractUserIDFromContext(ctx),
            UserTier:    extractUserTierFromContext(ctx),
        }
        
        if resp, err := s.textProvider.TranslateText(req); err == nil {
            translatedLines = append(translatedLines, resp.TranslatedText)
        } else {
            translatedLines = append(translatedLines, line)
        }
    }
    
    return strings.Join(translatedLines, "\n")
}

func extractContentFromMarkdown(line string) string {
    // إزالة علامات التنسيز
    line = strings.TrimSpace(line)
    
    // إزالة علامات العناوين
    for strings.HasPrefix(line, "#") {
        line = strings.TrimPrefix(line, "#")
    }
    
    // إزالة علامات القوائم
    if strings.HasPrefix(line, "* ") {
        line = strings.TrimPrefix(line, "* ")
    }
    if strings.HasPrefix(line, "- ") {
        line = strings.TrimPrefix(line, "- ")
    }
    
    // إزالة علامات الكود
    line = strings.Trim(line, "`")
    
    return strings.TrimSpace(line)
}

func applyTranslationToMarkdown(original string, originalContent string, translatedContent string) string {
    return strings.Replace(original, originalContent, translatedContent, 1)
}

func (s *TranslationService) translateHTML(ctx context.Context, html string, sourceLang string, targetLang string) string {
    // ترجمة بسيطة للـ HTML (الحفاظ على الوسوم)
    // في تطبيق حقيقي، تحتاج إلى parser HTML
    req := types.TranslationRequest{
        Text:        html,
        FromLang:    sourceLang,
        ToLang:      targetLang,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    resp, err := s.textProvider.TranslateText(req)
    if err != nil {
        return html
    }
    
    return resp.TranslatedText
}

func (s *TranslationService) translateJSON(ctx context.Context, jsonStr string, sourceLang string, targetLang string) string {
    // ترجمة قيم الـ JSON فقط، مع الحفاظ على المفاتيح
    // في تطبيق حقيقي، تحتاج إلى parse JSON
    req := types.TranslationRequest{
        Text:        jsonStr,
        FromLang:    sourceLang,
        ToLang:      targetLang,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    resp, err := s.textProvider.TranslateText(req)
    if err != nil {
        return jsonStr
    }
    
    return resp.TranslatedText
}

func (s *TranslationService) GetSupportedLanguages(ctx context.Context) ([]string, error) {
    // الحصول على اللغات المدعومة من المزود
    if provider, ok := s.textProvider.(interface{ GetSupportedLanguages() []string }); ok {
        return provider.GetSupportedLanguages(), nil
    }
    
    // قائمة افتراضية إذا لم يكن المزود يدعمها
    return []string{
        "en", "ar", "es", "fr", "de", "zh", "ja", "ko", "ru", "pt",
    }, nil
}

func (s *TranslationService) TranslateWithGlossary(ctx context.Context, text string, sourceLang string, targetLang string, glossary map[string]string) (*types.TranslationResponse, error) {
    // تطبيق المصطلحات قبل الترجمة
    for originalTerm, translation := range glossary {
        text = strings.ReplaceAll(text, originalTerm, fmt.Sprintf("[%s]", translation))
    }
    
    req := types.TranslationRequest{
        Text:        text,
        FromLang:    sourceLang,
        ToLang:      targetLang,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    resp, err := s.textProvider.TranslateText(req)
    if err != nil {
        return nil, err
    }
    
    // استبدال العلامات بالمصطلحات المترجمة
    translatedText := resp.TranslatedText
    for _, translation := range glossary {
        marker := fmt.Sprintf("[%s]", translation)
        translatedText = strings.ReplaceAll(translatedText, marker, translation)
    }
    
    return &types.TranslationResponse{
        TranslatedText: translatedText,
        Cost:           resp.Cost,
        Model:          resp.Model,
        CreatedAt:      resp.CreatedAt,
    }, nil
}

func (s *TranslationService) GetTranslationQuality(ctx context.Context, original string, translated string, sourceLang string, targetLang string) (float64, error) {
    // تقييم جودة الترجمة (بسيط)
    prompt := fmt.Sprintf(`Evaluate translation quality from %s to %s.
Original: %s
Translation: %s
Return a score from 0-100.`, sourceLang, targetLang, original, translated)
    
    textReq := types.TextRequest{
        Prompt:    prompt,
        MaxTokens: 10,
        UserID:    extractUserIDFromContext(ctx),
        UserTier:  extractUserTierFromContext(ctx),
    }
    
    resp, err := s.textProvider.GenerateText(textReq)
    if err != nil {
        return 0.0, fmt.Errorf("failed to evaluate translation quality: %w", err)
    }
    
    var score float64
    fmt.Sscanf(resp.Text, "%f", &score)
    
    return score, nil
}

func (s *TranslationService) GetServiceStats(ctx context.Context) map[string]interface{} {
    stats := make(map[string]interface{})
    
    // الحصول على إحصائيات من المزود إذا كانت متوفرة
    if provider, ok := s.textProvider.(interface{ GetStats() *types.ProviderStats }); ok {
        stats["provider_stats"] = provider.GetStats()
    }
    
    stats["service"] = "translation"
    stats["supported_languages"], _ = s.GetSupportedLanguages(ctx)
    stats["timestamp"] = time.Now()
    
    return stats
}