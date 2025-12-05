package services

import (
    "fmt"
    "strings"
)

type ContentService struct {
    textProvider TextProvider
}

func NewContentService(provider TextProvider) *ContentService {
    return &ContentService{textProvider: provider}
}

// GenerateBlogPost توليد مقال
func (s *ContentService) GenerateBlogPost(topic, language, tone string) (string, error) {
    prompt := fmt.Sprintf(`
    Write a comprehensive blog post about "%s"
    
    Requirements:
    - Language: %s
    - Tone: %s
    - Target audience: Entrepreneurs and business owners
    - Length: 1000-1500 words
    - Include headings, subheadings, and bullet points
    - Add a call-to-action at the end
    - Optimize for SEO with relevant keywords
    
    Structure:
    1. Introduction
    2. Main content with examples
    3. Actionable tips
    4. Conclusion
    `, topic, language, tone)
    
    return s.textProvider.GenerateText(prompt, "gemini-2.0-flash")
}

// GenerateSocialMediaPost توليد منشور وسائط اجتماعية
func (s *ContentService) GenerateSocialMediaPost(platform, topic, style string) (string, error) {
    platformPrompts := map[string]string{
        "linkedin": "professional, industry insights, business-focused",
        "twitter":  "concise, trending topics, include hashtags",
        "instagram": "engaging, visual descriptions, story-focused",
        "facebook": "community-oriented, conversational",
    }
    
    tone := platformPrompts[platform]
    if tone == "" {
        tone = "engaging"
    }
    
    prompt := fmt.Sprintf(`
    Create a %s post about "%s"
    
    Style: %s
    Platform: %s
    
    Include:
    - Main message
    - Supporting details
    - Call-to-action
    - Relevant hashtags (3-5)
    - Emojis if appropriate
    
    Make it shareable and engaging.
    `, platform, topic, style, platform)
    
    return s.textProvider.GenerateText(prompt, "llama3.2:3b")
}

// GenerateEmailCopy توليد نص بريد إلكتروني
func (s *ContentService) GenerateEmailCopy(purpose, audience, keyPoints string) (string, error) {
    prompt := fmt.Sprintf(`
    Write a professional email with the following:
    
    Purpose: %s
    Target Audience: %s
    Key Points to Include: %s
    
    Structure:
    - Subject line (compelling)
    - Greeting
    - Introduction
    - Main body
    - Call to action
    - Closing
    
    Tone: Professional yet friendly
    Length: 200-300 words
    `, purpose, audience, keyPoints)
    
    return s.textProvider.GenerateText(prompt, "mistral:7b")
}