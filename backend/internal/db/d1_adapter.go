package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/config"
)

// D1QueryRequest structure for proxy
type D1QueryRequest struct {
	SQL  string        `json:"sql"`
	Args []interface{} `json:"args,omitempty"`
}

type D1QueryResponse struct {
	Success bool                   `json:"success"`
	Result  []map[string]interface{} `json:"result,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// ExecuteD1ViaProxy sends SQL to a Cloudflare Worker proxy that executes against D1
func ExecuteD1ViaProxy(ctx context.Context, cfg *config.Config, sqlStatement string, args []interface{}) (*D1QueryResponse, error) {
	proxyURL := os.Getenv("D1_PROXY_URL")
	if proxyURL == "" {
		return nil, fmt.Errorf("D1_PROXY_URL is not set (you need a Worker proxy to run D1 queries from backend)")
	}

	reqBody := D1QueryRequest{
		SQL:  sqlStatement,
		Args: args,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", proxyURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	// إذا كان لدى البروكسي مفتاح سري:
	if secret := os.Getenv("D1_PROXY_SECRET"); secret != "" {
		req.Header.Set("X-Proxy-Secret", secret)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r D1QueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if !r.Success {
		return &r, fmt.Errorf("d1 proxy error: %s", r.Error)
	}
	return &r, nil
}