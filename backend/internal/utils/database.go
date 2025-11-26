package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/nawthtech/backend/internal/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	client *mongo.Client
	db     *mongo.Database
)

// InitDatabase تهيئة اتصال قاعدة البيانات
func InitDatabase(databaseURL string) error {
	if databaseURL == "" {
		logger.Stdout.Info("لا يوجد رابط قاعدة بيانات، العمل في وضع بدون قاعدة بيانات")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// تكوين خيارات العميل
	clientOptions := options.Client().
		ApplyURI(databaseURL).
		SetMaxPoolSize(100).                    // الحد الأقصى لعدد الاتصالات
		SetMinPoolSize(10).                     // الحد الأدنى لعدد الاتصالات
		SetMaxConnIdleTime(5 * time.Minute).    // أقصى وقت لبقاء الاتصال خاملاً
		SetConnectTimeout(10 * time.Second).    // وقت انتظار الاتصال
		SetServerSelectionTimeout(10 * time.Second) // وقت انتظار اختيار السيرفر

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("فشل في الاتصال بقاعدة البيانات: %w", err)
	}

	// التحقق من الاتصال
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("فشل في التحقق من اتصال قاعدة البيانات: %w", err)
	}

	db = client.Database("nawthtech")

	// إنشاء الفهارس
	if err := createIndexes(); err != nil {
		logger.Stderr.Warn("فشل في إنشاء بعض الفهارس", logger.ErrAttr(err))
	}

	logger.Stdout.Info("تم الاتصال بقاعدة البيانات بنجاح", 
		"database", "nawthtech",
		"max_pool_size", 100,
		"min_pool_size", 10,
	)

	return nil
}

// createIndexes إنشاء الفهارس الأساسية
func createIndexes() error {
	if db == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// فهارس للمستخدمين
	usersCollection := db.Collection("users")
	_, err := usersCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "phone", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	})
	if err != nil {
		return fmt.Errorf("فشل في إنشاء فهارس المستخدمين: %w", err)
	}

	// فهارس للخدمات
	servicesCollection := db.Collection("services")
	_, err = servicesCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "price", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "is_active", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	})
	if err != nil {
		return fmt.Errorf("فشل في إنشاء فهارس الخدمات: %w", err)
	}

	// فهارس للطلبات
	ordersCollection := db.Collection("orders")
	_, err = ordersCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "service_id", Value: 1}},
		},
	})
	if err != nil {
		return fmt.Errorf("فشل في إنشاء فهارس الطلبات: %w", err)
	}

	// فهارس للمدفوعات
	paymentsCollection := db.Collection("payments")
	_, err = paymentsCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "order_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	})
	if err != nil {
		return fmt.Errorf("فشل في إنشاء فهارس المدفوعات: %w", err)
	}

	logger.Stdout.Info("تم إنشاء الفهارس بنجاح")
	return nil
}

// GetDB الحصول على instance قاعدة البيانات
func GetDB() *mongo.Database {
	return db
}

// GetCollection الحصول على collection معين
func GetCollection(name string) *mongo.Collection {
	if db == nil {
		return nil
	}
	return db.Collection(name)
}

// WithTimeout إنشاء context مع timeout للعمليات
func WithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// IsConnected التحقق من اتصال قاعدة البيانات
func IsConnected() bool {
	if client == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := client.Ping(ctx, readpref.Primary())
	return err == nil
}

// GetDatabaseStats الحصول على إحصائيات قاعدة البيانات
func GetDatabaseStats() (bson.M, error) {
	if db == nil {
		return nil, fmt.Errorf("قاعدة البيانات غير متاحة")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result bson.M
	err := db.RunCommand(ctx, bson.D{{Key: "dbStats", Value: 1}}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("فشل في جلب إحصائيات قاعدة البيانات: %w", err)
	}

	return result, nil
}

// GetCollectionStats الحصول على إحصائيات collection معين
func GetCollectionStats(collectionName string) (bson.M, error) {
	if db == nil {
		return nil, fmt.Errorf("قاعدة البيانات غير متاحة")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection(collectionName)
	var result bson.M
	err := collection.FindOne(ctx, bson.D{{Key: "$collStats", Value: bson.M{"storageStats": bson.M{}}}}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("فشل في جلب إحصائيات الـ collection: %w", err)
	}

	return result, nil
}

// CloseDatabase إغلاق اتصالات قاعدة البيانات
func CloseDatabase() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		if err := client.Disconnect(ctx); err != nil {
			logger.Stderr.Error("فشل في إغلاق اتصال قاعدة البيانات", logger.ErrAttr(err))
		} else {
			logger.Stdout.Info("تم إغلاق اتصال قاعدة البيانات بنجاح")
		}
	}
}

// TransactionWrapper غلاف لتنفيذ العمليات ضمن transaction
func TransactionWrapper(callback func(sessionContext mongo.SessionContext) error) error {
	if client == nil {
		return fmt.Errorf("قاعدة البيانات غير متاحة")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	session, err := client.StartSession()
	if err != nil {
		return fmt.Errorf("فشل في بدء الجلسة: %w", err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessionContext mongo.SessionContext) (interface{}, error) {
		return nil, callback(sessionContext)
	})

	return err
}

// BulkOperations تنفيذ عمليات bulk
func BulkOperations(collectionName string, models []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	if db == nil {
		return nil, fmt.Errorf("قاعدة البيانات غير متاحة")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := db.Collection(collectionName)
	opts := options.BulkWrite().SetOrdered(false)
	
	result, err := collection.BulkWrite(ctx, models, opts)
	if err != nil {
		return nil, fmt.Errorf("فشل في تنفيذ عمليات الـ bulk: %w", err)
	}

	return result, nil
}

// CreateTextIndex إنشاء فهرس نصي للبحث
func CreateTextIndex(collectionName string, fields ...string) error {
	if db == nil {
		return fmt.Errorf("قاعدة البيانات غير متاحة")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.Collection(collectionName)
	
	// تحويل الحقول إلى نموذج فهرس نصي
	indexKeys := bson.D{}
	for _, field := range fields {
		indexKeys = append(indexKeys, bson.E{Key: field, Value: "text"})
	}

	indexModel := mongo.IndexModel{
		Keys:    indexKeys,
		Options: options.Index().SetName(fmt.Sprintf("%s_text_index", collectionName)),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("فشل في إنشاء الفهرس النصي: %w", err)
	}

	return nil
}

// HealthCheck فحص صحة اتصال قاعدة البيانات
func HealthCheck() map[string]interface{} {
	health := map[string]interface{}{
		"status":    "unhealthy",
		"timestamp": time.Now(),
	}

	if !IsConnected() {
		health["error"] = "قاعدة البيانات غير متصلة"
		return health
	}

	health["status"] = "healthy"

	// الحصول على إحصائيات إضافية
	if stats, err := GetDatabaseStats(); err == nil {
		health["stats"] = map[string]interface{}{
			"collections": stats["collections"],
			"objects":     stats["objects"],
			"dataSize":    stats["dataSize"],
			"storageSize": stats["storageSize"],
		}
	}

	return health
}
