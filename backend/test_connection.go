package backend

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	// الحصول على كلمة المرور من متغير البيئة
	mongoPassword := os.Getenv("MONGODB_PASSWORD")
	if mongoPassword == "" {
		log.Fatal("❌ MONGODB_PASSWORD environment variable is required")
	}

	// بناء URI مع كلمة المرور من المتغير
	uri := fmt.Sprintf("mongodb+srv://Nawthtech_db_user:%s@nawthtech-cluster.9nqbyeu.mongodb.net/nawthtech?retryWrites=true&w=majority&authSource=admin&authMechanism=SCRAM-SHA-256",
		mongoPassword)

	fmt.Println("🔗 Testing MongoDB Atlas connection...")
	fmt.Println("📡 URI:", maskURI(uri)) // إخفاء URI في الـ logs

	// إضافة Stable API
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPI)

	ctx, cancel := context.WithTimeout(context.Background(), 10)
	defer cancel()

	// الاتصال
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("⚠️ Failed to disconnect: %v", err)
		}
	}()

	// Ping للتحقق
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("❌ Failed to ping: %v", err)
	}

	fmt.Println("✅ Successfully connected to MongoDB Atlas!")
	fmt.Println("🏪 Database: nawthtech")
	fmt.Println("👤 User: Nawthtech_db_user")
}

// maskURI إخفاء كلمة المرور في URI للأمان
func maskURI(uri string) string {
	// إخفاء كلمة المرور في الـ logs
	const passwordPlaceholder = "***"
	start := "mongodb+srv://Nawthtech_db_user:"
	end := "@nawthtech-cluster"
	
	if len(uri) > len(start)+len(end) {
		passwordStart := len(start)
		passwordEnd := len(uri) - len(end)
		if passwordEnd > passwordStart {
			masked := uri[:passwordStart] + passwordPlaceholder + uri[passwordEnd:]
			return masked
		}
	}
	return "mongodb+srv://Nawthtech_db_user:***@nawthtech-cluster.9nqbyeu.mongodb.net/nawthtech..."
}