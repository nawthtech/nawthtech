package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// CloudflareEmailWorker handles email forwarding via Cloudflare Workers
type CloudflareEmailWorker struct {
	APIKey      string
	AccountID   string
	ScriptName  string
	ZoneID      string
	Domain      string
	AllowedList []string
	ForwardTo   string
}

// EmailMessage represents an incoming email
type EmailMessage struct {
	From    string            `json:"from"`
	To      []string          `json:"to"`
	Subject string            `json:"subject"`
	Text    string            `json:"text"`
	HTML    string            `json:"html"`
	Headers map[string]string `json:"headers"`
}

// CloudflareWorkerResponse represents the response from Cloudflare Workers API
type CloudflareWorkerResponse struct {
	Success bool        `json:"success"`
	Errors  []CFError   `json:"errors"`
	Result  interface{} `json:"result"`
}

type CFError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewCloudflareEmailWorker creates a new email worker instance
func NewCloudflareEmailWorker() (*CloudflareEmailWorker, error) {
	worker := &CloudflareEmailWorker{
		APIKey:      os.Getenv("CLOUDFLARE_API_KEY"),
		AccountID:   os.Getenv("CLOUDFLARE_ACCOUNT_ID"),
		ScriptName:  "nawthtech-email-router",
		ZoneID:      os.Getenv("CLOUDFLARE_ZONE_ID"),
		Domain:      "nawthtech.com",
		AllowedList: strings.Split(os.Getenv("EMAIL_ALLOWED_LIST"), ","),
		ForwardTo:   os.Getenv("EMAIL_FORWARD_TO"),
	}

	// Validate required environment variables
	if worker.APIKey == "" {
		return nil, fmt.Errorf("CLOUDFLARE_API_KEY is required")
	}
	if worker.AccountID == "" {
		return nil, fmt.Errorf("CLOUDFLARE_ACCOUNT_ID is required")
	}
	if worker.ForwardTo == "" {
		worker.ForwardTo = "inbox@corp" // Default fallback
	}

	return worker, nil
}

// DeployWorkerScript deploys the email worker to Cloudflare
func (w *CloudflareEmailWorker) DeployWorkerScript() error {
	workerScript := `export default {
  async email(message, env, ctx) {
    // Parse allowed list from environment
    const allowList = env.ALLOWED_LIST ? env.ALLOWED_LIST.split(',') : [];
    const forwardTo = env.FORWARD_TO || "inbox@corp";
    const domain = env.DOMAIN || "nawthtech.com";
    
    // Log incoming email for debugging
    console.log("Email received from:", message.from);
    console.log("Subject:", message.headers.get("subject"));
    
    // Check if sender is in allowed list
    const senderEmail = message.from.toLowerCase().trim();
    let isAllowed = false;
    
    for (const allowedEmail of allowList) {
      if (senderEmail === allowedEmail.toLowerCase().trim()) {
        isAllowed = true;
        break;
      }
    }
    
    // Also allow any email from the domain itself
    if (senderEmail.endsWith(domain)) {
      isAllowed = true;
    }
    
    if (!isAllowed) {
      console.log("Rejected email from:", message.from);
      message.setReject("Address not allowed");
      return;
    }
    
    try {
      // Forward to the specified address
      console.log("Forwarding email from", message.from, "to", forwardTo);
      await message.forward(forwardTo);
      
      // Send acknowledgment to sender (optional)
      if (env.SEND_ACKNOWLEDGMENT === "true") {
        await message.reply("Thank you for your email. It has been forwarded to the appropriate team.");
      }
      
    } catch (error) {
      console.error("Failed to forward email:", error);
      message.setReject("Failed to process email");
    }
  }
};`

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/workers/scripts/%s", w.AccountID, w.ScriptName)

	payload := map[string]interface{}{
		"script": workerScript,
		"bindings": []map[string]interface{}{
			{
				"type": "plain_text",
				"name": "ALLOWED_LIST",
				"text": strings.Join(w.AllowedList, ","),
			},
			{
				"type": "plain_text",
				"name": "FORWARD_TO",
				"text": w.ForwardTo,
			},
			{
				"type": "plain_text",
				"name": "DOMAIN",
				"text": w.Domain,
			},
			{
				"type": "plain_text",
				"name": "SEND_ACKNOWLEDGMENT",
				"text": "true",
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.APIKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to deploy worker: %v", err)
	}
	defer resp.Body.Close()

	var result CloudflareWorkerResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if !result.Success {
		return fmt.Errorf("failed to deploy worker: %v", result.Errors)
	}

	fmt.Println("✅ Email worker deployed successfully")
	return nil
}

// SetupDNSRecords sets up MX and TXT records for email routing
func (w *CloudflareEmailWorker) SetupDNSRecords() error {
	// 1. Create MX record
	mxRecord := map[string]interface{}{
		"type":     "MX",
		"name":     w.Domain,
		"content":  "route1.mx.cloudflare.net",
		"priority": 28,
		"ttl":      1, // Auto TTL
		"proxied":  false,
	}

	if err := w.createDNSRecord(mxRecord); err != nil {
		return fmt.Errorf("failed to create MX record: %v", err)
	}

	// 2. Create second MX record for redundancy
	mxRecord2 := map[string]interface{}{
		"type":     "MX",
		"name":     w.Domain,
		"content":  "route2.mx.cloudflare.net",
		"priority": 48,
		"ttl":      1,
		"proxied":  false,
	}

	if err := w.createDNSRecord(mxRecord2); err != nil {
		return fmt.Errorf("failed to create MX record 2: %v", err)
	}

	// 3. Create TXT record for SPF
	spfRecord := map[string]interface{}{
		"type":    "TXT",
		"name":    w.Domain,
		"content": "v=spf1 include:_spf.mx.cloudflare.net ~all",
		"ttl":     1,
		"proxied": false,
	}

	if err := w.createDNSRecord(spfRecord); err != nil {
		return fmt.Errorf("failed to create SPF record: %v", err)
	}

	// 4. Create DMARC record
	dmarcRecord := map[string]interface{}{
		"type":    "TXT",
		"name":    "_dmarc." + w.Domain,
		"content": "v=DMARC1; p=none; rua=mailto:dmarc-reports@nawthtech.com",
		"ttl":     1,
		"proxied": false,
	}

	if err := w.createDNSRecord(dmarcRecord); err != nil {
		return fmt.Errorf("failed to create DMARC record: %v", err)
	}

	fmt.Println("✅ DNS records for email routing configured successfully")
	return nil
}

func (w *CloudflareEmailWorker) createDNSRecord(record map[string]interface{}) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", w.ZoneID)

	jsonData, err := json.Marshal(record)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.APIKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result CloudflareWorkerResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if !result.Success {
		return fmt.Errorf("API error: %v", result.Errors)
	}

	return nil
}

// AddToAllowList adds an email to the allowed list
func (w *CloudflareEmailWorker) AddToAllowList(email string) error {
	// Check if email already exists
	for _, existing := range w.AllowedList {
		if strings.EqualFold(existing, email) {
			return fmt.Errorf("email %s already in allow list", email)
		}
	}

	w.AllowedList = append(w.AllowedList, email)

	// Update environment variable
	os.Setenv("EMAIL_ALLOWED_LIST", strings.Join(w.AllowedList, ","))

	// Redeploy worker with updated list
	return w.DeployWorkerScript()
}

// RemoveFromAllowList removes an email from the allowed list
func (w *CloudflareEmailWorker) RemoveFromAllowList(email string) error {
	var newList []string
	found := false

	for _, existing := range w.AllowedList {
		if !strings.EqualFold(existing, email) {
			newList = append(newList, existing)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("email %s not found in allow list", email)
	}

	w.AllowedList = newList
	os.Setenv("EMAIL_ALLOWED_LIST", strings.Join(w.AllowedList, ","))

	return w.DeployWorkerScript()
}

// GetAllowList returns the current allow list
func (w *CloudflareEmailWorker) GetAllowList() []string {
	return w.AllowedList
}

// TestEmailRouting sends a test email to verify setup
func (w *CloudflareEmailWorker) TestEmailRouting(testEmail string) error {
	// This would typically use a mail sending service
	// For now, just validate the setup
	if testEmail == "" {
		return fmt.Errorf("test email address is required")
	}

	fmt.Printf("Email routing test initialized\n")
	fmt.Printf("   Domain: %s\n", w.Domain)
	fmt.Printf("   Forward to: %s\n", w.ForwardTo)
	fmt.Printf("   Allowed senders: %v\n", w.AllowedList)
	fmt.Printf("   Test email: %s\n", testEmail)

	return nil
}