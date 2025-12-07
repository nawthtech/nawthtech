package services

import (
    "context"
    "fmt"
    "strings"
    "time"
    "github.com/nawthtech/nawthtech/backend/internal/ai/types"
)

type StrategyService struct {
    textProvider types.TextProvider
}

func NewStrategyService(provider types.TextProvider) *StrategyService {
    return &StrategyService{
        textProvider: provider,
    }
}

func (s *StrategyService) GenerateMarketingStrategy(ctx context.Context, product string, targetAudience string, budget string) (*types.TextResponse, error) {
    prompt := fmt.Sprintf(`Generate a comprehensive marketing strategy for: %s

Target Audience: %s
Budget: %s

Include the following sections:
1. Executive Summary
2. Target Market Analysis
3. Unique Selling Proposition
4. Marketing Channels
5. Content Strategy
6. Budget Allocation
7. Timeline
8. KPIs and Metrics

Make it detailed and actionable.`, product, targetAudience, budget)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   2000,
        Temperature: 0.7,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *StrategyService) GenerateBusinessPlan(ctx context.Context, businessIdea string, industry string) (*types.TextResponse, error) {
    prompt := fmt.Sprintf(`Generate a comprehensive business plan for: %s

Industry: %s

Include the following sections:
1. Executive Summary
2. Company Description
3. Market Analysis
4. Organization & Management
5. Product/Service Line
6. Marketing & Sales Strategy
7. Funding Request
8. Financial Projections
9. Appendix

Make it professional and investor-ready.`, businessIdea, industry)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   3000,
        Temperature: 0.7,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *StrategyService) GenerateContentStrategy(ctx context.Context, brand string, platforms []string, goals []string) (*types.TextResponse, error) {
    platformsStr := strings.Join(platforms, ", ")
    goalsStr := strings.Join(goals, ", ")
    
    prompt := fmt.Sprintf(`Generate a content strategy for: %s

Platforms: %s
Goals: %s

Include:
1. Content Pillars
2. Content Calendar (3 months)
3. Platform-Specific Strategies
4. Content Types and Formats
5. Distribution Plan
6. Engagement Tactics
7. Measurement and Analytics
8. Team and Resources

Make it practical and platform-optimized.`, brand, platformsStr, goalsStr)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   2500,
        Temperature: 0.7,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *StrategyService) GenerateSWOTAnalysis(ctx context.Context, company string, industry string) (*types.TextResponse, error) {
    prompt := fmt.Sprintf(`Generate a comprehensive SWOT analysis for: %s

Industry: %s

Structure the analysis as follows:

STRENGTHS (Internal, Positive Factors):
- List 5-7 strengths
- Include evidence and examples

WEAKNESSES (Internal, Negative Factors):
- List 5-7 weaknesses
- Be honest and constructive

OPPORTUNITIES (External, Positive Factors):
- List 5-7 opportunities
- Include market trends and gaps

THREATS (External, Negative Factors):
- List 5-7 threats
- Include competition and risks

STRATEGIC RECOMMENDATIONS:
- Based on the SWOT, provide 3-5 strategic recommendations
- Include short-term and long-term actions

Make it insightful and actionable.`, company, industry)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   1500,
        Temperature: 0.7,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *StrategyService) GenerateCompetitiveAnalysis(ctx context.Context, company string, competitors []string) (*types.TextResponse, error) {
    competitorsStr := strings.Join(competitors, ", ")
    
    prompt := fmt.Sprintf(`Generate a competitive analysis for: %s

Main Competitors: %s

Analyze each competitor in these areas:
1. Market Position and Share
2. Product/Service Offerings
3. Pricing Strategy
4. Marketing and Branding
5. Strengths and Weaknesses
6. Customer Reviews and Sentiment
7. Recent Developments

Then provide:
1. Competitive Matrix (Comparison Table)
2. Competitive Advantage Recommendations
3. Market Opportunity Identification
4. Strategic Positioning Advice

Make it data-driven and strategic.`, company, competitorsStr)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   2000,
        Temperature: 0.7,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *StrategyService) GenerateProductLaunchPlan(ctx context.Context, product string, targetMarket string, launchDate string) (*types.TextResponse, error) {
    prompt := fmt.Sprintf(`Generate a product launch plan for: %s

Target Market: %s
Launch Date: %s

Include:
1. Pre-Launch Activities (90-30 days before)
   - Market Research
   - Beta Testing
   - Building Hype
   - Influencer Outreach

2. Launch Phase (30-0 days before)
   - Marketing Campaign
   - Press Release
   - Social Media Blitz
   - Launch Event

3. Post-Launch Activities (0-90 days after)
   - Customer Support
   - Feedback Collection
   - Performance Analysis
   - Iteration Planning

4. Budget and Resources
5. Risk Management
6. Success Metrics

Make it detailed with specific timelines and responsibilities.`, product, targetMarket, launchDate)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   2500,
        Temperature: 0.7,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *StrategyService) GenerateBrandPositioning(ctx context.Context, brand string, targetCustomer string, competitors []string) (*types.TextResponse, error) {
    competitorsStr := strings.Join(competitors, ", ")
    
    prompt := fmt.Sprintf(`Generate brand positioning for: %s

Target Customer: %s
Competitors: %s

Include:
1. Brand Essence (Core Identity)
2. Brand Promise (Value Proposition)
3. Brand Personality
4. Target Audience Personas
5. Competitive Differentiation
6. Positioning Statement
7. Key Messages
8. Tone of Voice Guidelines
9. Visual Identity Direction
10. Implementation Roadmap

Make it distinctive and memorable.`, brand, targetCustomer, competitorsStr)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   2000,
        Temperature: 0.7,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *StrategyService) GenerateCrisisManagementPlan(ctx context.Context, organization string, potentialCrises []string) (*types.TextResponse, error) {
    crisesStr := strings.Join(potentialCrises, ", ")
    
    prompt := fmt.Sprintf(`Generate a crisis management plan for: %s

Potential Crises to Address: %s

Include:
1. Crisis Management Team Structure
2. Communication Protocols
3. Immediate Response Checklist
4. Stakeholder Communication Templates
5. Media Relations Strategy
6. Social Media Response Guidelines
7. Recovery and Reputation Management
8. Post-Crisis Evaluation
9. Training and Preparedness
10. Legal and Compliance Considerations

Make it comprehensive and ready-to-use.`, organization, crisesStr)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   2000,
        Temperature: 0.7,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *StrategyService) GenerateDigitalTransformationStrategy(ctx context.Context, company string, currentState string, goals []string) (*types.TextResponse, error) {
    goalsStr := strings.Join(goals, ", ")
    
    prompt := fmt.Sprintf(`Generate a digital transformation strategy for: %s

Current State: %s
Transformation Goals: %s

Include:
1. Vision and Objectives
2. Technology Assessment
3. Change Management Plan
4. Implementation Roadmap (Phased Approach)
5. Digital Skills Development
6. Data and Analytics Strategy
7. Cybersecurity Considerations
8. Budget and ROI Analysis
9. Success Metrics
10. Risk Management

Make it practical with clear milestones.`, company, currentState, goalsStr)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   2500,
        Temperature: 0.7,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *StrategyService) GenerateOKRs(ctx context.Context, department string, timeframe string, companyGoals []string) (*types.TextResponse, error) {
    companyGoalsStr := strings.Join(companyGoals, ", ")
    
    prompt := fmt.Sprintf(`Generate Objectives and Key Results (OKRs) for: %s

Timeframe: %s
Company Goals to Align With: %s

Structure as:
OBJECTIVE 1: [Clear, inspirational objective]
- Key Result 1.1: [Measurable outcome]
- Key Result 1.2: [Measurable outcome]
- Key Result 1.3: [Measurable outcome]

OBJECTIVE 2: [Clear, inspirational objective]
- Key Result 2.1: [Measurable outcome]
- Key Result 2.2: [Measurable outcome]
- Key Result 2.3: [Measurable outcome]

OBJECTIVE 3: [Clear, inspirational objective]
- Key Result 3.1: [Measurable outcome]
- Key Result 3.2: [Measurable outcome]
- Key Result 3.3: [Measurable outcome]

Include:
1. How to measure success
2. Initiatives to achieve each KR
3. Dependencies and risks
4. Review cadence

Make the OKRs SMART (Specific, Measurable, Achievable, Relevant, Time-bound).`, department, timeframe, companyGoalsStr)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   1500,
        Temperature: 0.7,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *StrategyService) GetServiceStats(ctx context.Context) map[string]interface{} {
    stats := make(map[string]interface{})
    
    // الحصول على إحصائيات من المزود إذا كانت متوفرة
    if provider, ok := s.textProvider.(interface{ GetStats() *types.ProviderStats }); ok {
        stats["provider_stats"] = provider.GetStats()
    }
    
    stats["service"] = "strategy"
    stats["capabilities"] = []string{
        "marketing_strategy",
        "business_plan",
        "content_strategy", 
        "swot_analysis",
        "competitive_analysis",
        "product_launch_plan",
        "brand_positioning",
        "crisis_management",
        "digital_transformation",
        "okrs",
    }
    stats["timestamp"] = time.Now()
    
    return stats
}