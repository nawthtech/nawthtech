package mongodb

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nawthtech/nawthtech/backend/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ================================
// هياكل البيانات
// ================================

// MongoDBConfig إعدادات MongoDB
type MongoDBConfig struct {
	URI               string
	DatabaseName      string
	TestDatabaseName  string
	ConnectTimeout    time.Duration
	OperationTimeout  time.Duration
	MaxPoolSize       uint64
	MinPoolSize       uint64
	SocketTimeout     time.Duration
	ServerSelectionTimeout time.Duration
	ReplicaSet        string
	SSL               bool
}

// MongoDBService خدمة MongoDB
type MongoDBService struct {
	Client   *mongo.Client
	Database *mongo.Database
	Config   *MongoDBConfig
	mu       sync.RWMutex
}

// CollectionStats إحصائيات المجموعة
type CollectionStats struct {
	Name        string `bson:"ns" json:"name"`
	Count       int64  `bson:"count" json:"count"`
	Size        int64  `bson:"size" json:"size"`
	AvgObjSize  int64  `bson:"avgObjSize" json:"avg_obj_size"`
	StorageSize int64  `bson:"storageSize" json:"storage_size"`
	TotalIndexSize int64 `bson:"totalIndexSize" json:"total_index_size"`
}

// DatabaseStats إحصائيات قاعدة البيانات
type DatabaseStats struct {
	DB          string    `bson:"db" json:"db"`
	Collections int64     `bson:"collections" json:"collections"`
	Objects     int64     `bson:"objects" json:"objects"`
	DataSize    int64     `bson:"dataSize" json:"data_size"`
	StorageSize int64     `bson:"storageSize" json:"storage_size"`
	IndexSize   int64     `bson:"indexSize" json:"index_size"`
	IndexCount  int64     `bson:"indexes" json:"index_count"`
}

// QueryResult نتيجة الاستعلام
type QueryResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Count   int64       `json:"count,omitempty"`
}

// ================================
// دوال التهيئة
// ================================

// NewMongoDBConfig إنشاء إعدادات MongoDB جديدة
func NewMongoDBConfig() *MongoDBConfig {
	return &MongoDBConfig{
		URI:               getEnv("MONGODB_URI", "mongodb://localhost:27017/nawthtech"),
		DatabaseName:      getEnv("MONGODB_DB_NAME", "nawthtech"),
		TestDatabaseName:  getEnv("MONGODB_TEST_DB_NAME", "nawthtech_test"),
		ConnectTimeout:    time.Duration(getEnvInt("MONGODB_CONNECT_TIMEOUT", 10)) * time.Second,
		OperationTimeout:  time.Duration(getEnvInt("MONGODB_OPERATION_TIMEOUT", 30)) * time.Second,
		MaxPoolSize:       uint64(getEnvInt("MONGODB_MAX_POOL_SIZE", 100)),
		MinPoolSize:       uint64(getEnvInt("MONGODB_MIN_POOL_SIZE", 10)),
		SocketTimeout:     time.Duration(getEnvInt("MONGODB_SOCKET_TIMEOUT", 30)) * time.Second,
		ServerSelectionTimeout: time.Duration(getEnvInt("MONGODB_SERVER_SELECTION_TIMEOUT", 10)) * time.Second,
		SSL:               getEnvBool("MONGODB_SSL", false),
	}
}

// NewMongoDBService إنشاء خدمة MongoDB جديدة
func NewMongoDBService() (*MongoDBService, error) {
	config := NewMongoDBConfig()
	
	// استخدام قاعدة بيانات الاختبار في بيئة التطوير
	if os.Getenv("APP_ENV") == "test" {
		config.DatabaseName = config.TestDatabaseName
	}

	client, database, err := connectMongoDB(config)
	if err != nil {
		return nil, fmt.Errorf("فشل في الاتصال بـ MongoDB: %v", err)
	}

	service := &MongoDBService{
		Client:   client,
		Database: database,
		Config:   config,
	}

	logger.Info(context.Background(), "✅ تم تهيئة خدمة MongoDB بنجاح",
		"database", config.DatabaseName,
		"environment", getEnv("APP_ENV", "development"),
		"connection_string", maskConnectionString(config.URI),
	)

	return service, nil
}

// connectMongoDB الاتصال بقاعدة بيانات MongoDB
func connectMongoDB(config *MongoDBConfig) (*mongo.Client, *mongo.Database, error) {
	startTime := time.Now()

	// إعداد خيارات العميل
	clientOptions := options.Client().
		ApplyURI(config.URI).
		SetConnectTimeout(config.ConnectTimeout).
		SetSocketTimeout(config.SocketTimeout).
		SetServerSelectionTimeout(config.ServerSelectionTimeout).
		SetMaxPoolSize(config.MaxPoolSize).
		SetMinPoolSize(config.MinPoolSize).
		SetReadPreference(readpref.Primary()).
		SetRetryWrites(true)

	// إعداد SSL إذا كان مفعلاً
	if config.SSL {
		clientOptions.SetTLSConfig(nil) // يمكن تخصيص إعدادات SSL هنا
	}

	// الاتصال بقاعدة البيانات
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("فشل في إنشاء اتصال MongoDB: %v", err)
	}

	// اختبار الاتصال
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, nil, fmt.Errorf("فشل في اختبار اتصال MongoDB: %v", err)
	}

	database := client.Database(config.DatabaseName)

	logger.Info(context.Background(), "✅ تم الاتصال بـ MongoDB بنجاح",
		"duration", time.Since(startTime),
		"database", config.DatabaseName,
	)

	return client, database, nil
}

// ================================
// دوال العمليات الأساسية
// ================================

// GetCollection الحصول على مجموعة
func (s *MongoDBService) GetCollection(name string) *mongo.Collection {
	return s.Database.Collection(name)
}

// InsertOne إدراج وثيقة واحدة
func (s *MongoDBService) InsertOne(ctx context.Context, collectionName string, document interface{}) (*mongo.InsertOneResult, error) {
	startTime := time.Now()
	collection := s.GetCollection(collectionName)

	ctx, cancel := context.WithTimeout(ctx, s.Config.OperationTimeout)
	defer cancel()

	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		logger.Error(ctx, "❌ فشل في إدراج وثيقة",
			"collection", collectionName,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return nil, err
	}

	logger.Debug(ctx, "✅ تم إدراج وثيقة بنجاح",
		"collection", collectionName,
		"duration", time.Since(startTime),
		"inserted_id", result.InsertedID,
	)

	return result, nil
}

// FindOne البحث عن وثيقة واحدة
func (s *MongoDBService) FindOne(ctx context.Context, collectionName string, filter interface{}, result interface{}) error {
	startTime := time.Now()
	collection := s.GetCollection(collectionName)

	ctx, cancel := context.WithTimeout(ctx, s.Config.OperationTimeout)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logger.Debug(ctx, "⚠️ لم يتم العثور على وثيقة",
				"collection", collectionName,
				"duration", time.Since(startTime),
			)
			return err
		}
		
		logger.Error(ctx, "❌ فشل في البحث عن وثيقة",
			"collection", collectionName,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return err
	}

	logger.Debug(ctx, "✅ تم العثور على وثيقة بنجاح",
		"collection", collectionName,
		"duration", time.Since(startTime),
	)

	return nil
}

// Find البحث عن عدة وثائق
func (s *MongoDBService) Find(ctx context.Context, collectionName string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	startTime := time.Now()
	collection := s.GetCollection(collectionName)

	ctx, cancel := context.WithTimeout(ctx, s.Config.OperationTimeout)
	defer cancel()

	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		logger.Error(ctx, "❌ فشل في البحث عن وثائق",
			"collection", collectionName,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return nil, err
	}

	logger.Debug(ctx, "✅ تم البحث عن وثائق بنجاح",
		"collection", collectionName,
		"duration", time.Since(startTime),
	)

	return cursor, nil
}

// UpdateOne تحديث وثيقة واحدة
func (s *MongoDBService) UpdateOne(ctx context.Context, collectionName string, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	startTime := time.Now()
	collection := s.GetCollection(collectionName)

	ctx, cancel := context.WithTimeout(ctx, s.Config.OperationTimeout)
	defer cancel()

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error(ctx, "❌ فشل في تحديث وثيقة",
			"collection", collectionName,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return nil, err
	}

	logger.Debug(ctx, "✅ تم تحديث وثيقة بنجاح",
		"collection", collectionName,
		"duration", time.Since(startTime),
		"matched_count", result.MatchedCount,
		"modified_count", result.ModifiedCount,
	)

	return result, nil
}

// DeleteOne حذف وثيقة واحدة
func (s *MongoDBService) DeleteOne(ctx context.Context, collectionName string, filter interface{}) (*mongo.DeleteResult, error) {
	startTime := time.Now()
	collection := s.GetCollection(collectionName)

	ctx, cancel := context.WithTimeout(ctx, s.Config.OperationTimeout)
	defer cancel()

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		logger.Error(ctx, "❌ فشل في حذف وثيقة",
			"collection", collectionName,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return nil, err
	}

	logger.Debug(ctx, "✅ تم حذف وثيقة بنجاح",
		"collection", collectionName,
		"duration", time.Since(startTime),
		"deleted_count", result.DeletedCount,
	)

	return result, nil
}

// Count حساب عدد الوثائق
func (s *MongoDBService) Count(ctx context.Context, collectionName string, filter interface{}) (int64, error) {
	startTime := time.Now()
	collection := s.GetCollection(collectionName)

	ctx, cancel := context.WithTimeout(ctx, s.Config.OperationTimeout)
	defer cancel()

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error(ctx, "❌ فشل في حساب الوثائق",
			"collection", collectionName,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return 0, err
	}

	logger.Debug(ctx, "✅ تم حساب الوثائق بنجاح",
		"collection", collectionName,
		"duration", time.Since(startTime),
		"count", count,
	)

	return count, nil
}

// ================================
// دوال متقدمة
// ================================

// Aggregate تنفيذ عملية تجميع
func (s *MongoDBService) Aggregate(ctx context.Context, collectionName string, pipeline interface{}) (*mongo.Cursor, error) {
	startTime := time.Now()
	collection := s.GetCollection(collectionName)

	ctx, cancel := context.WithTimeout(ctx, s.Config.OperationTimeout*2) // وقت أطول للتجميع
	defer cancel()

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Error(ctx, "❌ فشل في تنفيذ عملية التجميع",
			"collection", collectionName,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return nil, err
	}

	logger.Debug(ctx, "✅ تم تنفيذ عملية التجميع بنجاح",
		"collection", collectionName,
		"duration", time.Since(startTime),
	)

	return cursor, nil
}

// BulkWrite تنفيذ عمليات مجمعة
func (s *MongoDBService) BulkWrite(ctx context.Context, collectionName string, operations []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	startTime := time.Now()
	collection := s.GetCollection(collectionName)

	ctx, cancel := context.WithTimeout(ctx, s.Config.OperationTimeout*2)
	defer cancel()

	result, err := collection.BulkWrite(ctx, operations)
	if err != nil {
		logger.Error(ctx, "❌ فشل في تنفيذ العمليات المجمعة",
			"collection", collectionName,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return nil, err
	}

	logger.Debug(ctx, "✅ تم تنفيذ العمليات المجمعة بنجاح",
		"collection", collectionName,
		"duration", time.Since(startTime),
		"operations_count", len(operations),
		"inserted_count", result.InsertedCount,
		"modified_count", result.ModifiedCount,
		"deleted_count", result.DeletedCount,
	)

	return result, nil
}

// CreateIndex إنشاء فهرس
func (s *MongoDBService) CreateIndex(ctx context.Context, collectionName string, keys interface{}, opts ...*options.IndexOptions) (string, error) {
	startTime := time.Now()
	collection := s.GetCollection(collectionName)

	indexModel := mongo.IndexModel{
		Keys:    keys,
		Options: options.MergeIndexOptions(opts...),
	}

	ctx, cancel := context.WithTimeout(ctx, s.Config.OperationTimeout)
	defer cancel()

	indexName, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		logger.Error(ctx, "❌ فشل في إنشاء فهرس",
			"collection", collectionName,
			"duration", time.Since(startTime),
			"error", err.Error(),
		)
		return "", err
	}

	logger.Info(ctx, "✅ تم إنشاء فهرس بنجاح",
		"collection", collectionName,
		"duration", time.Since(startTime),
		"index_name", indexName,
	)

	return indexName, nil
}

// ================================
// دوال الإدارة والمراقبة
// ================================

// HealthCheck فحص صحة اتصال MongoDB
func (s *MongoDBService) HealthCheck(ctx context.Context) map[string]interface{} {
	startTime := time.Now()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// فحص الاتصال
	err := s.Client.Ping(ctx, readpref.Primary())
	if err != nil {
		return map[string]interface{}{
			"service":   "mongodb",
			"status":    "error",
			"timestamp": time.Now().UTC(),
			"error":     err.Error(),
			"duration":  time.Since(startTime).String(),
		}
	}

	// الحصول على إحصائيات قاعدة البيانات
	stats, err := s.GetDatabaseStats(ctx)
	if err != nil {
		return map[string]interface{}{
			"service":   "mongodb",
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"warning":   "فشل في الحصول على الإحصائيات: " + err.Error(),
			"duration":  time.Since(startTime).String(),
		}
	}

	return map[string]interface{}{
		"service":   "mongodb",
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"database":  s.Config.DatabaseName,
		"stats":     stats,
		"duration":  time.Since(startTime).String(),
	}
}

// GetDatabaseStats الحصول على إحصائيات قاعدة البيانات
func (s *MongoDBService) GetDatabaseStats(ctx context.Context) (*DatabaseStats, error) {
	var stats DatabaseStats
	command := bson.D{{Key: "dbStats", Value: 1}}

	err := s.Database.RunCommand(ctx, command).Decode(&stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetCollectionStats الحصول على إحصائيات المجموعة
func (s *MongoDBService) GetCollectionStats(ctx context.Context, collectionName string) (*CollectionStats, error) {
	var stats CollectionStats
	command := bson.D{{Key: "collStats", Value: collectionName}}

	err := s.Database.RunCommand(ctx, command).Decode(&stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// ListCollections سرد جميع المجموعات
func (s *MongoDBService) ListCollections(ctx context.Context) ([]string, error) {
	var collections []string

	cursor, err := s.Database.ListCollections(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result struct {
			Name string `bson:"name"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		collections = append(collections, result.Name)
	}

	return collections, nil
}

// ================================
// دوال المساعدة
// ================================

// Close إغلاق اتصال MongoDB
func (s *MongoDBService) Close() error {
	if s.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.Client.Disconnect(ctx)
		if err != nil {
			logger.Error(context.Background(), "❌ فشل في إغلاق اتصال MongoDB", "error", err.Error())
			return err
		}

		logger.Info(context.Background(), "✅ تم إغلاق اتصال MongoDB بنجاح")
	}
	return nil
}

// GetClient الحصول على عميل MongoDB
func (s *MongoDBService) GetClient() *mongo.Client {
	return s.Client
}

// GetDatabase الحصول على قاعدة البيانات
func (s *MongoDBService) GetDatabase() *mongo.Database {
	return s.Database
}

// ObjectIDFromString تحويل سلسلة إلى ObjectID
func ObjectIDFromString(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

// ObjectIDToString تحويل ObjectID إلى سلسلة
func ObjectIDToString(id primitive.ObjectID) string {
	return id.Hex()
}

// ================================
// دوال مساعدة للبيئة
// ================================

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// maskConnectionString إخفاء كلمة السر في رابط الاتصال للأمان
func maskConnectionString(connectionString string) string {
	if len(connectionString) > 50 {
		return connectionString[:30] + "****" + connectionString[len(connectionString)-20:]
	}
	return "****"
}