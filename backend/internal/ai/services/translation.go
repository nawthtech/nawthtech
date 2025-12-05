package services

import (
    "fmt"
    "strings"
)

type TranslationService struct {
    textProvider TextProvider
}

func NewTranslationService(provider TextProvider) *TranslationService {
    return &TranslationService{textProvider: provider}
}

// TranslateText ترجمة نص
func (s *TranslationService) TranslateText(text, sourceLang, targetLang string) (string, error) {
    prompt := fmt.Sprintf(`
    Translate the following text from %s to %s:
    
    "%s"
    
    Translation requirements:
    - Maintain original meaning and tone
    - Use natural, fluent language
    - Preserve technical terms when appropriate
    - Adapt cultural references if needed
    - Keep formatting (headings, lists, etc.)
    
    Return only the translated text.
    `, sourceLang, targetLang, text)
    
    return s.textProvider.GenerateText(prompt, "qwen2.5:7b")
}

// LocalizeContent توطين المحتوى
func (s *TranslationService) LocalizeContent(content, targetCulture string) (string, error) {
    prompt := fmt.Sprintf(`
    Localize this content for %s culture:
    
    "%s"
    
    Localization requirements:
    - Translate if needed
    - Adapt cultural references
    - Use appropriate idioms
    - Adjust humor and tone
    - Consider local customs and sensitivities
    - Format dates, numbers, and currencies correctly
    
    Make it feel native to the target culture.
    `, targetCulture, content)
    
    return s.textProvider.GenerateText(prompt, "gemini-2.0-flash")
}