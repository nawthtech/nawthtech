package services

import (
    "fmt"
)

type AnalysisService struct {
    textProvider TextProvider
}

func NewAnalysisService(provider TextProvider) *AnalysisService {
    return &AnalysisService{textProvider: provider}
}

// AnalyzeMarketTrends تحليل اتجاهات السوق
func (s *AnalysisService) AnalyzeMarketTrends(industry, timeframe string) (string, error) {
    prompt := fmt.Sprintf(`
    Analyze market trends for the %s industry over %s.
    
    Provide insights on:
    1. Current market size and growth rate
    2. Key drivers and challenges
    3. Emerging technologies
    4. Competitive landscape
    5. Future predictions
    6. Recommendations for businesses
    
    Format as a structured report with clear sections.
    `, industry, timeframe)
    
    return s.textProvider.GenerateText(prompt, "qwen2.5:7b")
}

// AnalyzeAudience تحليل الجمهور المستهدف
func (s *AnalysisService) AnalyzeAudience(demographics, interests, behavior string) (string, error) {
    prompt := fmt.Sprintf(`
    Analyze the target audience with these characteristics:
    
    Demographics: %s
    Interests: %s
    Online Behavior: %s
    
    Provide insights on:
    - Content preferences
    - Platform usage
    - Purchase triggers
    - Pain points
    - Communication style
    - Personalization opportunities
    
    Make it actionable for digital marketing.
    `, demographics, interests, behavior)
    
    return s.textProvider.GenerateText(prompt, "gemini-2.0-flash")
}

// SWOTAnalysis تحليل SWOT
func (s *AnalysisService) SWOTAnalysis(businessType, marketPosition string) (string, error) {
    prompt := fmt.Sprintf(`
    Conduct a comprehensive SWOT analysis for a %s business.
    
    Market Position: %s
    
    Structure:
    
    STRENGTHS:
    - List 5-7 internal strengths
    
    WEAKNESSES:
    - List 5-7 internal weaknesses
    
    OPPORTUNITIES:
    - List 5-7 external opportunities
    
    THREATS:
    - List 5-7 external threats
    
    RECOMMENDATIONS:
    - Strategic recommendations based on analysis
    
    Be specific and actionable.
    `, businessType, marketPosition)
    
    return s.textProvider.GenerateText(prompt, "gemini-2.0-flash")
}