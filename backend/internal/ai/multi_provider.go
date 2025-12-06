// backend/internal/ai/multi_provider.go
package ai

type ProviderType string

const (
    ProviderGemini ProviderType = "gemini"
    ProviderOpenAI ProviderType = "openai"
    ProviderOllama ProviderType = "ollama"
)

type MultiProvider struct {
    providers map[ProviderType]AIProvider
    strategy  RoutingStrategy
}

type RoutingStrategy interface {
    SelectProvider(userTier string, promptType string) ProviderType
}

func NewMultiProvider() *MultiProvider {
    return &MultiProvider{
        providers: map[ProviderType]AIProvider{
            ProviderGemini: NewGeminiProvider(),
            ProviderOpenAI: NewOpenAIProvider(),
            ProviderOllama: NewOllamaProvider(),
        },
        strategy: &TieredStrategy{},
    }
}

func (m *MultiProvider) Generate(userID, prompt string) (string, error) {
    userTier := getUserTier(userID)
    promptType := classifyPrompt(prompt)
    
    providerType := m.strategy.SelectProvider(userTier, promptType)
    provider, exists := m.providers[providerType]
    
    if !exists || !provider.IsAvailable() {
        // Fallback to Ollama
        provider = m.providers[ProviderOllama]
    }
    
    return provider.Generate(prompt)
}

// إستراتيجية التوجيه حسب خطة المستخدم
type TieredStrategy struct{}

func (s *TieredStrategy) SelectProvider(userTier, promptType string) ProviderType {
    switch userTier {
    case "free":
        return ProviderGemini  // مجاني للمستخدمين المجانيين
    case "premium":
        if promptType == "analysis" || promptType == "strategy" {
            return ProviderOpenAI  // GPT للأمور المتقدمة
        }
        return ProviderGemini
    case "enterprise":
        return ProviderOpenAI  // الأفضل للشركات
    default:
        return ProviderGemini
    }
}