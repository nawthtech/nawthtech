// backend/cmd/email-setup/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	
	"github.com/nawthtech/nawthtech/backend/internal/email"
)

func main() {
	action := flag.String("action", "setup", "setup, test, add-email, list")
	emailAddr := flag.String("email", "", "Email address for add/test")
	flag.Parse()

	// Initialize email service
	service, err := email.NewService()
	if err != nil {
		log.Fatalf("❌ Failed to initialize email service: %v", err)
	}

	switch *action {
	case "setup":
		fmt.Println("🚀 Setting up email routing...")
		if err := service.SetupEmailRouting(); err != nil {
			log.Fatalf("❌ Setup failed: %v", err)
		}
	case "test":
		if *emailAddr == "" {
			log.Fatal("❌ Email address required for test")
		}
		fmt.Printf("🧪 Testing email: %s\n", *emailAddr)
		if err := service.TestEmail(*emailAddr); err != nil {
			log.Fatalf("❌ Test failed: %v", err)
		}
	case "add-email":
		if *emailAddr == "" {
			log.Fatal("❌ Email address required")
		}
		if err := service.AddToAllowList(*emailAddr); err != nil {
			log.Fatalf("❌ Failed to add email: %v", err)
		}
	case "list":
		emails := service.GetAllowList()
		fmt.Println("📋 Allowed emails:")
		for _, e := range emails {
			fmt.Printf("  • %s\n", e)
		}
	default:
		fmt.Println("Available actions: setup, test, add-email, list")
		os.Exit(1)
	}
	
	fmt.Println("✅ Done!")
}