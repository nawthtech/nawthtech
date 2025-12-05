// backend/internal/ai/video/svd.go
package video

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
)

type SVDClient struct {
    apiKey string
    baseURL string
}

func NewSVDClient() *SVDClient {
    return &SVDClient{
        apiKey:  os.Getenv("STABILITY_API_KEY"), // مجاني محدود
        baseURL: "https://api.stability.ai/v2alpha/generation",
    }
}

func (c *SVDClient) GenerateVideo(prompt string, duration int) ([]byte, error) {
    reqBody := map[string]interface{}{
        "text_prompts": []map[string]interface{}{
            {
                "text": prompt,
                "weight": 1.0,
            },
        },
        "cfg_scale": 7,
        "steps":     50,
        "seed":      0,
    }
    
    jsonBody, _ := json.Marshal(reqBody)
    
    req, _ := http.NewRequest("POST", c.baseURL+"/video-to-video", 
        bytes.NewBuffer(jsonBody))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    return io.ReadAll(resp.Body)
}

// مجاني: 25 توليد/شهر (مع API key)
// بدون API: استخدام محلي