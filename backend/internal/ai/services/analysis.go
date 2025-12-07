package services

import (
    "context"
    "fmt"
    "strings"
    "time"
    "github.com/nawthtech/nawthtech/backend/internal/ai/types"
)

type AnalysisService struct {
    textProvider types.TextProvider
}

func NewAnalysisService(provider types.TextProvider) *AnalysisService {
    return &AnalysisService{
        textProvider: provider,
    }
}

// AnalyzeMarketTrends تحليل اتجاهات السوق
func (s *AnalysisService) AnalyzeMarketTrends(ctx context.Context, industry string, timeframe string) (*types.TextResponse, error) {
    prompt := fmt.Sprintf(`Analyze market trends for the %s industry over the %s timeframe.

Provide a comprehensive analysis including:

1. CURRENT MARKET SITUATION:
   - Market size and growth rate
   - Key players and market share
   - Current challenges and opportunities

2. TREND ANALYSIS:
   - Emerging technologies
   - Consumer behavior shifts
   - Regulatory changes
   - Economic factors

3. COMPETITIVE LANDSCAPE:
   - Major competitors analysis
   - New entrants and disruptors
   - Competitive strategies

4. GROWTH OPPORTUNITIES:
   - Untapped market segments
   - Geographic opportunities
   - Product/service innovations

5. RISK ASSESSMENT:
   - Market risks
   - Competitive threats
   - Economic/political risks

6. STRATEGIC RECOMMENDATIONS:
   - Short-term actions (next 3-6 months)
   - Medium-term strategies (6-18 months)
   - Long-term positioning (18+ months)

7. KEY METRICS TO WATCH:
   - Industry-specific KPIs
   - Leading indicators
   - Lagging indicators

Provide data-driven insights with specific examples and actionable recommendations.`, industry, timeframe)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   2000,
        Temperature: 0.6,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *AnalysisService) AnalyzeSentiment(ctx context.Context, text string, sourceType string) (*types.AnalysisResponse, error) {
    prompt := fmt.Sprintf(`Perform sentiment analysis on the following %s text:

Text: %s

Analyze:
1. Overall sentiment (Positive/Negative/Neutral)
2. Sentiment intensity (Score 0-100)
3. Key emotional indicators
4. Specific positive aspects mentioned
5. Specific negative aspects mentioned
6. Neutral or mixed signals
7. Tone analysis (Formal/Informal, Emotional/Rational, etc.)
8. Potential biases or loaded language

Provide the analysis in a structured format with confidence scores.`, sourceType, text)
    
    req := types.AnalysisRequest{
        Text:     text,
        Prompt:   prompt,
        UserID:   extractUserIDFromContext(ctx),
        UserTier: extractUserTierFromContext(ctx),
    }
    
    // نحتاج إلى معرفة إذا كان المزود يدعم AnalyzeText مباشرة
    if provider, ok := s.textProvider.(interface{ AnalyzeText(req types.AnalysisRequest) (*types.AnalysisResponse, error) }); ok {
        return provider.AnalyzeText(req)
    }
    
    // خيار احتياطي: استخدام GenerateText
    textReq := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   500,
        Temperature: 0.3,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    resp, err := s.textProvider.GenerateText(textReq)
    if err != nil {
        return nil, err
    }
    
    // تحويل النتيجة إلى AnalysisResponse
    return &types.AnalysisResponse{
        Result:     resp.Text,
        Confidence: 0.85, // قيمة افتراضية
        Cost:       resp.Cost,
        Model:      resp.ModelUsed,
        CreatedAt:  resp.CreatedAt,
    }, nil
}

func (s *AnalysisService) AnalyzeCompetitor(ctx context.Context, competitor string, yourCompany string, analysisType string) (*types.TextResponse, error) {
    prompt := fmt.Sprintf(`Perform a comprehensive competitor analysis for %s compared to %s.

Analysis type: %s

Include:

1. COMPETITOR OVERVIEW:
   - Company background
   - Market position
   - Core offerings
   - Target audience

2. STRENGTHS ANALYSIS:
   - Product/service strengths
   - Brand strengths
   - Operational strengths
   - Financial strengths

3. WEAKNESSES ANALYSIS:
   - Product/service gaps
   - Brand weaknesses
   - Operational challenges
   - Financial vulnerabilities

4. OPPORTUNITIES (for your company):
   - Competitor weaknesses to exploit
   - Market gaps they're missing
   - Customer pain points they're not addressing
   - Technological opportunities

5. THREATS (from this competitor):
   - Competitive advantages they have
   - Market share they could take
   - Strategic moves they might make
   - Pricing pressure

6. COMPETITIVE ADVANTAGE MATRIX:
   - Compare across key dimensions
   - Visual ranking (1-5 scale)

7. STRATEGIC RECOMMENDATIONS:
   - Defensive strategies
   - Offensive strategies
   - Partnership opportunities
   - Innovation priorities

Provide actionable intelligence.`, competitor, yourCompany, analysisType)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   2500,
        Temperature: 0.6,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *AnalysisService) AnalyzeFinancialData(ctx context.Context, financialMetrics map[string]interface{}, timeframe string) (*types.TextResponse, error) {
    metricsStr := ""
    for key, value := range financialMetrics {
        metricsStr += fmt.Sprintf("- %s: %v\n", key, value)
    }
    
    prompt := fmt.Sprintf(`Analyze the following financial data:

Timeframe: %s

Financial Metrics:
%s

Provide analysis covering:

1. FINANCIAL HEALTH ASSESSMENT:
   - Liquidity analysis
   - Solvency assessment
   - Profitability analysis
   - Efficiency metrics

2. TREND ANALYSIS:
   - Revenue growth trends
   - Profit margin trends
   - Cost structure analysis
   - Cash flow patterns

3. BENCHMARK COMPARISON:
   - Industry benchmarks
   - Historical performance
   - Competitor comparisons

4. RISK ANALYSIS:
   - Financial risks identified
   - Sustainability concerns
   - Dependence analysis

5. STRENGTHS AND WEAKNESSES:
   - Financial strengths
   - Areas for improvement
   - Red flags

6. RECOMMENDATIONS:
   - Immediate actions
   - Strategic improvements
   - Investment priorities
   - Risk mitigation

7. FORECASTING:
   - Short-term projections
   - Scenario analysis
   - Sensitivity analysis

Provide data-driven insights with specific recommendations.`, timeframe, metricsStr)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   2000,
        Temperature: 0.5,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *AnalysisService) AnalyzeCustomerFeedback(ctx context.Context, feedback []string, product string) (*types.AnalysisResponse, error) {
    feedbackStr := strings.Join(feedback, "\n")
    
    prompt := fmt.Sprintf(`Analyze customer feedback for: %s

Customer Feedback:
%s

Analyze:

1. OVERALL SATISFACTION:
   - Overall sentiment score
   - Satisfaction trends
   - Key drivers of satisfaction

2. THEMATIC ANALYSIS:
   - Recurring themes (positive)
   - Recurring issues (negative)
   - Feature requests
   - Usability feedback

3. PRIORITIZATION:
   - Critical issues (urgent action needed)
   - Important improvements (short-term)
   - Enhancement requests (long-term)
   - Low-priority feedback

4. CUSTOMER SEGMENTATION:
   - Power users vs casual users
   - Different use cases
   - Demographic patterns

5. COMPETITIVE INSIGHTS:
   - Mentions of competitors
   - Comparative advantages
   - Market positioning feedback

6. ACTIONABLE INSIGHTS:
   - Product improvements
   - Service enhancements
   - Communication improvements
   - Training needs

7. ROI CALCULATION:
   - Impact of addressing top issues
   - Cost-benefit analysis
   - Prioritization matrix

Provide structured analysis with specific action items.`, product, feedbackStr)
    
    req := types.AnalysisRequest{
        Text:     feedbackStr,
        Prompt:   prompt,
        UserID:   extractUserIDFromContext(ctx),
        UserTier: extractUserTierFromContext(ctx),
    }
    
    if provider, ok := s.textProvider.(interface{ AnalyzeText(req types.AnalysisRequest) (*types.AnalysisResponse, error) }); ok {
        return provider.AnalyzeText(req)
    }
    
    textReq := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   1500,
        Temperature: 0.4,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    resp, err := s.textProvider.GenerateText(textReq)
    if err != nil {
        return nil, err
    }
    
    return &types.AnalysisResponse{
        Result:     resp.Text,
        Confidence: 0.9,
        Cost:       resp.Cost,
        Model:      resp.ModelUsed,
        CreatedAt:  resp.CreatedAt,
        Categories: []string{"customer_feedback", "sentiment", "product_analysis"},
    }, nil
}

func (s *AnalysisService) AnalyzeSWOT(ctx context.Context, organization string, industry string) (*types.TextResponse, error) {
    prompt := fmt.Sprintf(`Perform a comprehensive SWOT analysis for %s in the %s industry.

Provide detailed analysis in this structure:

STRENGTHS (Internal Positive Factors):
- List 5-8 key strengths
- Evidence and examples for each
- Competitive advantages derived
- Sustainability of each strength

WEAKNESSES (Internal Negative Factors):
- List 5-8 key weaknesses
- Honest assessment of limitations
- Impact on competitiveness
- Potential for improvement

OPPORTUNITIES (External Positive Factors):
- List 5-8 key opportunities
- Market trends enabling opportunities
- Competitive gaps
- Growth potential assessment

THREATS (External Negative Factors):
- List 5-8 key threats
- Competitive threats
- Market threats
- Regulatory/economic threats

CROSS-ANALYSIS:
- Strength-Opportunity Strategies (SO)
- Strength-Threat Strategies (ST)
- Weakness-Opportunity Strategies (WO)
- Weakness-Threat Strategies (WT)

STRATEGIC PRIORITIZATION:
- High-impact, high-feasibility actions
- Quick wins
- Long-term strategic moves
- Risk mitigation strategies

IMPLEMENTATION ROADMAP:
- Phase 1 (0-3 months)
- Phase 2 (3-12 months)
- Phase 3 (12+ months)

Make it actionable with specific recommendations.`, organization, industry)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   3000,
        Temperature: 0.6,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *AnalysisService) AnalyzeRisk(ctx context.Context, project string, riskAreas []string) (*types.TextResponse, error) {
    riskAreasStr := strings.Join(riskAreas, ", ")
    
    prompt := fmt.Sprintf(`Perform a comprehensive risk analysis for: %s

Risk areas to consider: %s

Analysis Structure:

1. RISK IDENTIFICATION:
   - List all potential risks (categorized)
   - Likelihood assessment for each
   - Impact assessment for each
   - Risk score calculation (Likelihood × Impact)

2. RISK CATEGORIZATION:
   - Strategic risks
   - Operational risks
   - Financial risks
   - Compliance risks
   - Reputational risks
   - Technological risks

3. RISK PRIORITIZATION:
   - High priority risks (immediate action)
   - Medium priority risks (monitor closely)
   - Low priority risks (accept/monitor)

4. RISK MITIGATION STRATEGIES:
   - Prevention strategies
   - Reduction strategies
   - Transfer strategies (insurance, outsourcing)
   - Acceptance strategies

5. CONTINGENCY PLANNING:
   - "What if" scenarios
   - Emergency response plans
   - Business continuity plans

6. RISK MONITORING:
   - Key risk indicators (KRIs)
   - Monitoring frequency
   - Reporting structure
   - Escalation procedures

7. RISK APPETITE AND TOLERANCE:
   - Organization's risk appetite
   - Risk tolerance levels
   - Risk-adjusted return expectations

8. RECOMMENDATIONS:
   - Immediate actions
   - Long-term risk management improvements
   - Resource allocation recommendations

Provide a practical, actionable risk management framework.`, project, riskAreasStr)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   2500,
        Temperature: 0.5,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *AnalysisService) AnalyzePerformanceMetrics(ctx context.Context, metrics map[string]interface{}, benchmarks map[string]interface{}) (*types.TextResponse, error) {
    metricsStr := ""
    for key, value := range metrics {
        metricsStr += fmt.Sprintf("- %s: %v\n", key, value)
    }
    
    benchmarksStr := ""
    for key, value := range benchmarks {
        benchmarksStr += fmt.Sprintf("- %s: %v\n", key, value)
    }
    
    prompt := fmt.Sprintf(`Analyze performance metrics against benchmarks:

Performance Metrics:
%s

Benchmarks:
%s

Analysis Requirements:

1. PERFORMANCE ASSESSMENT:
   - Overall performance rating
   - Areas exceeding benchmarks
   - Areas below benchmarks
   - Performance trends

2. GAP ANALYSIS:
   - Significant gaps identified
   - Root cause analysis for gaps
   - Impact of each gap

3. STRENGTHS IDENTIFICATION:
   - Top performing areas
   - Competitive advantages
   - Best practices identified

4. IMPROVEMENT OPPORTUNITIES:
   - Priority improvement areas
   - Quick win opportunities
   - Long-term improvement areas

5. BENCHMARK VALIDATION:
   - Benchmark relevance assessment
   - Industry comparison validity
   - Suggested benchmark adjustments

6. CORRELATION ANALYSIS:
   - Relationship between different metrics
   - Leading vs lagging indicators
   - Cause-effect relationships

7. PREDICTIVE ANALYSIS:
   - Future performance projections
   - Scenario analysis
   - Target setting recommendations

8. ACTION PLAN:
   - Immediate corrective actions
   - Process improvements
   - Resource allocation recommendations
   - Timeline for improvements

Provide data-driven insights with specific recommendations.`, metricsStr, benchmarksStr)
    
    req := types.TextRequest{
        Prompt:      prompt,
        MaxTokens:   2000,
        Temperature: 0.5,
        UserID:      extractUserIDFromContext(ctx),
        UserTier:    extractUserTierFromContext(ctx),
    }
    
    return s.textProvider.GenerateText(req)
}

func (s *AnalysisService) GetServiceStats(ctx context.Context) map[string]interface{} {
    stats := make(map[string]interface{})
    
    // الحصول على إحصائيات من المزود إذا كانت متوفرة
    if provider, ok := s.textProvider.(interface{ GetStats() *types.ProviderStats }); ok {
        stats["provider_stats"] = provider.GetStats()
    }
    
    stats["service"] = "analysis"
    stats["analysis_types"] = []string{
        "market_trends",
        "sentiment",
        "competitor", 
        "financial",
        "customer_feedback",
        "swot",
        "risk",
        "performance_metrics",
    }
    stats["timestamp"] = time.Now()
    
    return stats
}