package models

// AudioModel نموذج صوتي
type AudioModel struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Provider    string   `json:"provider"`
    MaxDuration int      `json:"max_duration"`
    Languages   []string `json:"languages"`
    CostPerMinute float64 `json:"cost_per_minute"`
    IsLocal     bool     `json:"is_local"`
}

// النماذج المجانية للصوت
var FreeAudioModels = []AudioModel{
    {
        ID:          "bark",
        Name:        "Bark",
        Provider:    "Suno AI",
        MaxDuration: 14,
        Languages:   []string{"en", "ar", "es", "fr", "de"},
        CostPerMinute: 0.0,
        IsLocal:     true,
    },
    {
        ID:          "xtts",
        Name:        "XTTS-v2",
        Provider:    "Coqui AI",
        MaxDuration: 0, // غير محدود
        Languages:   []string{"en", "ar", "es", "fr", "de", "zh"},
        CostPerMinute: 0.0,
        IsLocal:     true,
    },
    {
        ID:          "whisper",
        Name:        "Whisper",
        Provider:    "OpenAI",
        MaxDuration: 0,
        Languages:   []string{"multilingual"},
        CostPerMinute: 0.0,
        IsLocal:     true,
    },
    {
        ID:          "mms",
        Name:        "MMS (Massively Multilingual Speech)",
        Provider:    "Meta",
        MaxDuration: 0,
        Languages:   []string{"ar", "en", "fr", "es", "zh", "ru", "pt"},
        CostPerMinute: 0.0,
        IsLocal:     true,
    },
}