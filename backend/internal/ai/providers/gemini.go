package providers

import (
    "context"
    "fmt"
    "os"
    "strings"
    
    "google.golang.org/genai"
)

type GeminiProvider struct {
    client *genai.Client
}

func NewGeminiProvider() (*GeminiProvider, error) {
    apiKey := os.Getenv("GEMINI_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("GEMINI_API_KEY not set")
    }
    
    ctx := context.Background()
    client, err := genai.NewClient(ctx, &genai.ClientConfig{
        APIKey: apiKey,
    })
    if err != nil {
        return nil, err
    }
    
    return &GeminiProvider{client: client}, nil
}

func (p *GeminiProvider) GenerateText(prompt string, model string) (string, error) {
    ctx := context.Background()
    
    if model == "" {
        model = "gemini-2.0-flash"
    }
    
    result, err := p.client.Models.GenerateContent(ctx, model, genai.Text(prompt), nil)
    if err != nil {
        return "", err
    }
    
    return extractText(result), nil
}

func extractText(result *genai.GenerateContentResponse) string {
    var texts []string
    for _, candidate := range result.Candidates {
        for _, part := range candidate.Content.Parts {
            if part.Text != "" {
                texts = append(texts, part.Text)
            }
        }
    }
    return strings.Join(texts, "\n")
}

// مجاني: 60 request/دقيقة
// 1M tokens/شهر