// backend/internal/email/service.go
package email

import (
	"fmt"
	"os"
	"strings"
)

type Service struct {
	config *Config
}

type Config struct {
	CloudflareAPIKey string
	CloudflareZoneID string
	Domain           string
	ForwardTo        string
	AllowedList      []string
}

func NewService() (*Service, error) {
	config := &Config{
		CloudflareAPIKey: os.Getenv("CLOUDFLARE_API_KEY"),
		CloudflareZoneID: os.Getenv("CLOUDFLARE_ZONE_ID"),
		Domain:           "nawthtech.com",
		ForwardTo:        os.Getenv("EMAIL_FORWARD_TO"),
		AllowedList:      strings.Split(os.Getenv("EMAIL_ALLOWED_LIST"), ","),
	}

	if config.CloudflareAPIKey == "" || config.CloudflareZoneID == "" {
		return nil, fmt.Errorf("Cloudflare credentials not set")
	}

	return &Service{config: config}, nil
}

func (s *Service) SetupEmailRouting() error {
	// Setup DNS records
	if err := s.setupDNS(); err != nil {
		return fmt.Errorf("DNS setup failed: %v", err)
	}

	// Setup Cloudflare Email Routing
	if err := s.setupCloudflareRouting(); err != nil {
		return fmt.Errorf("Cloudflare routing setup failed: %v", err)
	}

	return nil
}

func (s *Service) setupDNS() error {
	// Implementation for DNS setup
	fmt.Println("🌐 Setting up DNS records...")
	// ... DNS API calls
	return nil
}

func (s *Service) setupCloudflareRouting() error {
	fmt.Println("📧 Configuring Cloudflare Email Routing...")
	// ... Cloudflare API calls
	return nil
}

func (s *Service) TestEmail(email string) error {
	fmt.Printf("Sending test email to %s...\n", email)
	// Implementation
	return nil
}

func (s *Service) AddToAllowList(email string) error {
	s.config.AllowedList = append(s.config.AllowedList, email)
	// Update environment/config
	return nil
}

func (s *Service) GetAllowList() []string {
	return s.config.AllowedList
}