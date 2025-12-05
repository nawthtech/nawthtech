package video

import (
    "fmt"
    "os"
    "sync"
)

type HybridVideoProvider struct {
    providers []VideoProvider
    mu        sync.RWMutex
}

func NewHybridVideoProvider() *HybridVideoProvider {
    h := &HybridVideoProvider{}
    
    // 1. أولاً: محاولة النماذج المحلية (مجانية)
    if os.Getenv("ENABLE_LOCAL_VIDEO") == "true" {
        if svd := NewLocalSVD(); svd != nil {
            h.providers = append(h.providers, svd)
            fmt.Println("✅ Local SVD provider initialized")
        }
    }
    
    // 2. ثانياً: خدمات مجانية محدودة
    if apiKey := os.Getenv("STABILITY_API_KEY"); apiKey != "" {
        h.providers = append(h.providers, NewStabilityClient())
        fmt.Println("✅ Stability AI provider initialized")
    }
    
    // 3. أخيراً: Veo المدفوع (للجودة العالية)
    if os.Getenv("ENABLE_VEO") == "true" {
        if veo, err := NewVeoProvider(); err == nil {
            h.providers = append(h.providers, veo)
            fmt.Println("✅ Google Veo provider initialized")
        }
    }
    
    return h
}

func (h *HybridVideoProvider) GenerateVideo(req VideoRequest) (*VideoResponse, error) {
    // 1. محاولة المحلي أولاً (مجاني)
    for _, provider := range h.providers {
        if provider.IsLocal() && provider.IsAvailable() {
            fmt.Println("🎬 Using local video generation (free)")
            return provider.GenerateVideo(req)
        }
    }
    
    // 2. محاولة خدمات مجانية محدودة
    for _, provider := range h.providers {
        if provider.IsFree() && provider.IsAvailable() {
            fmt.Println("🎬 Using free cloud video generation")
            return provider.GenerateVideo(req)
        }
    }
    
    // 3. استخدام المدفوع إذا فشل الجميع
    for _, provider := range h.providers {
        if provider.IsAvailable() {
            fmt.Println("🎬 Using paid video generation")
            return provider.GenerateVideo(req)
        }
    }
    
    return nil, fmt.Errorf("no video providers available")
}