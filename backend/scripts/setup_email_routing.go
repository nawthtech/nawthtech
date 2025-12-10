package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type CloudflareEmailConfig struct {
	APIKey    string
	ZoneID    string
	Domain    string
	DestEmail string
}

func main() {
	config := CloudflareEmailConfig{
		APIKey:    os.Getenv("CLOUDFLARE_API_KEY"),
		ZoneID:    os.Getenv("CLOUDFLARE_ZONE_ID"),
		Domain:    "nawthtech.com",
		DestEmail: os.Getenv("DESTINATION_EMAIL"),
	}

	if config.APIKey == "" || config.ZoneID == "" {
		log.Fatal("❌ Cloudflare API credentials are required")
	}

	fmt.Println("🚀 Setting up Cloudflare Email Routing for nawthtech.com")

	// 1. Setup DNS Records
	if err := setupDNSRecords(config); err != nil {
		log.Fatalf("❌ Failed to setup DNS: %v", err)
	}

	fmt.Println("✅ DNS records configured")

	// 2. Create Custom Address
	if err := createCustomAddress(config); err != nil {
		log.Fatalf("❌ Failed to create custom address: %v", err)
	}

	fmt.Println("✅ Custom address created")
	fmt.Println("📧 Please check your email to confirm the routing")
}

func setupDNSRecords(config CloudflareEmailConfig) error {
	// MX Records
	mxRecords := []map[string]interface{}{
		{
			"type":     "MX",
			"name":     config.Domain,
			"content":  "route1.mx.cloudflare.net",
			"priority": 28,
			"ttl":      1,
			"proxied":  false,
		},
		{
			"type":     "MX",
			"name":     config.Domain,
			"content":  "route2.mx.cloudflare.net",
			"priority": 48,
			"ttl":      1,
			"proxied":  false,
		},
	}

	// TXT Records
	txtRecords := []map[string]interface{}{
		{
			"type":    "TXT",
			"name":    config.Domain,
			"content": "v=spf1 include:_spf.mx.cloudflare.net ~all",
			"ttl":     1,
			"proxied": false,
		},
		{
			"type":    "TXT",
			"name":    "_dmarc." + config.Domain,
			"content": "v=DMARC1; p=none; rua=mailto:dmarc-reports@" + config.Domain,
			"ttl":     1,
			"proxied": false,
		},
	}

	// Create all records
	allRecords := append(mxRecords, txtRecords...)

	for _, record := range allRecords {
		if err := createDNSRecord(config, record); err != nil {
			return fmt.Errorf("failed to create record %s: %v", record["type"], err)
		}
		fmt.Printf("  ✅ Created %s record\n", record["type"])
	}

	return nil
}

func createDNSRecord(config CloudflareEmailConfig, record map[string]interface{}) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", config.ZoneID)

	jsonData, err := json.Marshal(record)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	return nil
}

func createCustomAddress(config CloudflareEmailConfig) error {
	// Note: Email Routing API is currently in beta
	// This endpoint might change
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/email/routing/addresses", config.ZoneID)

	addressData := map[string]interface{}{
		"email":    fmt.Sprintf("*@%s", config.Domain),
		"verified": false,
	}

	jsonData, err := json.Marshal(addressData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create custom address: %d", resp.StatusCode)
	}

	return nil
}
