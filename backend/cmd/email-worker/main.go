package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	
	"nawthtech/internal/email"
)

func main() {
	// Parse command line flags
	action := flag.String("action", "deploy", "Action to perform: deploy, setup-dns, add-email, remove-email, list, test")
	emailAddr := flag.String("email", "", "Email address for add/remove/test actions")
	flag.Parse()

	// Initialize email worker
	worker, err := email.NewCloudflareEmailWorker()
	if err != nil {
		log.Fatalf("❌ Failed to initialize email worker: %v", err)
	}

	// Execute requested action
	switch *action {
	case "deploy":
		fmt.Println("🚀 Deploying email worker...")
		if err := worker.DeployWorkerScript(); err != nil {
			log.Fatalf("❌ Failed to deploy worker: %v", err)
		}
		
	case "setup-dns":
		fmt.Println("🌐 Setting up DNS records for email...")
		if err := worker.SetupDNSRecords(); err != nil {
			log.Fatalf("❌ Failed to setup DNS: %v", err)
		}
		
	case "add-email":
		if *emailAddr == "" {
			log.Fatal("❌ Email address is required for add-email action")
		}
		fmt.Printf("➕ Adding %s to allow list...\n", *emailAddr)
		if err := worker.AddToAllowList(*emailAddr); err != nil {
			log.Fatalf("❌ Failed to add email: %v", err)
		}
		fmt.Println("✅ Email added successfully")
		
	case "remove-email":
		if *emailAddr == "" {
			log.Fatal("❌ Email address is required for remove-email action")
		}
		fmt.Printf("➖ Removing %s from allow list...\n", *emailAddr)
		if err := worker.RemoveFromAllowList(*emailAddr); err != nil {
			log.Fatalf("❌ Failed to remove email: %v", err)
		}
		fmt.Println("✅ Email removed successfully")
		
	case "list":
		fmt.Println("📋 Current allow list:")
		emails := worker.GetAllowList()
		if len(emails) == 0 {
			fmt.Println("   No emails in allow list")
		} else {
			for i, email := range emails {
				fmt.Printf("   %d. %s\n", i+1, email)
			}
		}
		
	case "test":
		if *emailAddr == "" {
			log.Fatal("❌ Test email address is required for test action")
		}
		fmt.Printf("🧪 Testing email routing for %s...\n", *emailAddr)
		if err := worker.TestEmailRouting(*emailAddr); err != nil {
			log.Fatalf("❌ Test failed: %v", err)
		}
		fmt.Println("✅ Test completed successfully")
		
	default:
		fmt.Println("❌ Unknown action. Available actions:")
		fmt.Println("   deploy      - Deploy email worker script")
		fmt.Println("   setup-dns   - Setup DNS records for email routing")
		fmt.Println("   add-email   - Add email to allow list")
		fmt.Println("   remove-email - Remove email from allow list")
		fmt.Println("   list        - List allowed emails")
		fmt.Println("   test        - Test email routing")
		os.Exit(1)
	}
	
	fmt.Println("🎉 Operation completed successfully!")
}