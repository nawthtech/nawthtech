package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nawthtech/backend/internal/config"
	"github.com/nawthtech/backend/internal/logger"
)

type WorkerService struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	enabled    bool
}

type WorkerRequest struct {
	Method  string      `json:"method"`
	Path    string      `json:"path"`
	Headers http.Header `json:"headers"`
	Body    interface{} `json:"body,omitempty"`
}

type WorkerResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func NewWorkerService(cfg *config.Config) (*WorkerService, error) {
	workerURL := os.Getenv("WORKER_API_URL")
	if workerURL == "" {
		workerURL = "https://api.nawthtech.com"
	}

	apiKey := os.Getenv("WORKER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("WORKER_API_KEY is required for worker service")
	}

	enabled := os.Getenv("WORKER_ENABLED") != "false"

	return &WorkerService{
		baseURL: workerURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		enabled: enabled,
	}, nil
}

func (ws *WorkerService) Call(ctx context.Context, method, path string, data interface{}) (*WorkerResponse, error) {
	if !ws.enabled {
		return nil, fmt.Errorf("worker service is disabled")
	}

	url := fmt.Sprintf("%s%s", ws.baseURL, path)
	
	var reqBody []byte
	if data != nil {
		var err error
		reqBody, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+ws.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Requested-With", "Go-Backend")

	logger.Debug(ctx, "calling worker API", 
		"method", method, 
		"path", path, 
		"url", url,
	)

	resp, err := ws.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("worker API call failed: %v", err)
	}
	defer resp.Body.Close()

	var workerResp WorkerResponse
	if err := json.NewDecoder(resp.Body).Decode(&workerResp); err != nil {
		return nil, fmt.Errorf("failed to decode worker response: %v", err)
	}

	if resp.StatusCode >= 400 {
		return &workerResp, fmt.Errorf("worker API error: %s", workerResp.Error)
	}

	return &workerResp, nil
}

func (ws *WorkerService) HealthCheck() (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := ws.Call(ctx, "GET", "/health", nil)
	if err != nil {
		return nil, err
	}

	healthData, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid health response format")
	}

	return healthData, nil
}

// User operations
func (ws *WorkerService) CreateUser(ctx context.Context, userData map[string]interface{}) (*WorkerResponse, error) {
	return ws.Call(ctx, "POST", "/api/v1/auth/register", userData)
}

func (ws *WorkerService) LoginUser(ctx context.Context, credentials map[string]string) (*WorkerResponse, error) {
	return ws.Call(ctx, "POST", "/api/v1/auth/login", credentials)
}

func (ws *WorkerService) GetCurrentUser(ctx context.Context, token string) (*WorkerResponse, error) {
	// Note: Token should be passed in headers, not body
	return ws.Call(ctx, "GET", "/api/v1/auth/me", nil)
}

func (ws *WorkerService) GenerateAI(ctx context.Context, prompt string, options map[string]interface{}) (*WorkerResponse, error) {
	request := map[string]interface{}{
		"prompt":  prompt,
		"options": options,
	}
	return ws.Call(ctx, "POST", "/api/v1/ai/generate", request)
}

func (ws *WorkerService) CreateService(ctx context.Context, serviceData map[string]interface{}) (*WorkerResponse, error) {
	return ws.Call(ctx, "POST", "/api/v1/services", serviceData)
}

func (ws *WorkerService) GetServices(ctx context.Context, params map[string]string) (*WorkerResponse, error) {
	// Build query string
	path := "/api/v1/services"
	if len(params) > 0 {
		path += "?"
		for key, value := range params {
			path += fmt.Sprintf("%s=%s&", key, value)
		}
		path = path[:len(path)-1] // Remove last &
	}
	return ws.Call(ctx, "GET", path, nil)
}