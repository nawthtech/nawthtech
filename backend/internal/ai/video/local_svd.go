// backend/internal/ai/video/local_svd.go
package video

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "image"
    "image/jpeg"
    "io"
    "net/http"
    "os"
    "time"
)

type LocalSVDProvider struct {
    apiURL string
    client *http.Client
}

func NewLocalSVDProvider() *LocalSVDProvider {
    apiURL := os.Getenv("SVD_API_URL")
    if apiURL == "" {
        apiURL = "http://localhost:7860"
    }
    
    return &LocalSVDProvider{
        apiURL: apiURL,
        client: &http.Client{Timeout: 300 * time.Second},
    }
}

func (p *LocalSVDProvider) GenerateVideoFromImage(imgData []byte, prompt string) ([]byte, error) {
    // تحويل الصورة إلى base64
    imgBase64 := base64.StdEncoding.EncodeToString(imgData)
    
    reqBody := map[string]interface{}{
        "image": imgBase64,
        "prompt": prompt,
        "num_frames": 14,
        "fps": 7,
        "seed": -1,
    }
    
    jsonBody, _ := json.Marshal(reqBody)
    
    resp, err := p.client.Post(p.apiURL+"/api/predict", 
        "application/json", bytes.NewBuffer(jsonBody))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result struct {
        Data []string `json:"data"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    if len(result.Data) == 0 {
        return nil, fmt.Errorf("no video generated")
    }
    
    // تحويل base64 إلى bytes
    videoData, err := base64.StdEncoding.DecodeString(result.Data[0])
    if err != nil {
        return nil, err
    }
    
    return videoData, nil
}

func (p *LocalSVDProvider) IsAvailable() bool {
    resp, err := p.client.Get(p.apiURL + "/health")
    return err == nil && resp.StatusCode == 200
}

func (p *LocalSVDProvider) IsLocal() bool {
    return true
}

func (p *LocalSVDProvider) IsFree() bool {
    return true
}