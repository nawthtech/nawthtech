package services

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestNewServiceContainer اختبار إنشاء حاوية الخدمات
func TestNewServiceContainer(t *testing.T) {
	// Test with nil database (for unit tests)
	container := NewServiceContainer(nil)
	assert.NotNil(t, container, "Service container should be created")

	// Test that all services are initialized
	assert.NotNil(t, container.Auth, "Auth service should be initialized")
	assert.NotNil(t, container.User, "User service should be initialized")
	assert.NotNil(t, container.Service, "Service service should be initialized")
	assert.NotNil(t, container.Category, "Category service should be initialized")
	assert.NotNil(t, container.Order, "Order service should be initialized")
	assert.NotNil(t, container.Payment, "Payment service should be initialized")
	assert.NotNil(t, container.Upload, "Upload service should be initialized")
	assert.NotNil(t, container.Notification, "Notification service should be initialized")
	assert.NotNil(t, container.Admin, "Admin service should be initialized")
	assert.NotNil(t, container.Cache, "Cache service should be initialized")
}

// TestNewServiceContainerWithConfig اختبار إنشاء حاوية الخدمات مع الإعدادات
func TestNewServiceContainerWithConfig(t *testing.T) {
	// Create a mock config
	config := struct {
		Environment string
		Version     string
	}{
		Environment: "test",
		Version:     "1.0.0",
	}

	// Create a mock logger
	logger := struct {
		Info  func(msg string, fields ...interface{})
		Error func(msg string, fields ...interface{})
	}{
		Info:  func(msg string, fields ...interface{}) {},
		Error: func(msg string, fields ...interface{}) {},
	}

	// Test with nil database
	container := NewServiceContainerWithConfig(nil, config, logger)
	assert.NotNil(t, container, "Service container should be created with config")

	// Test db field
	assert.Nil(t, container.db, "Database should be nil when passing nil")
}

// TestCacheServiceMethods اختبار طرق خدمة التخزين المؤقت
func TestCacheServiceMethods(t *testing.T) {
	// Create cache service
	cacheService := NewCacheService()
	assert.NotNil(t, cacheService, "Cache service should be created")

	// Test basic cache operations
	testKey := "test_key"
	testValue := "test_value"
	expiration := 1 * time.Minute

	// Test Set method
	err := cacheService.Set(testKey, testValue, expiration)
	assert.NoError(t, err, "Set should not return error")

	// Test Get method
	value, err := cacheService.Get(testKey)
	assert.NoError(t, err, "Get should not return error")
	assert.Equal(t, testValue, value, "Retrieved value should match stored value")

	// Test Exists method
	exists, err := cacheService.Exists(testKey)
	assert.NoError(t, err, "Exists should not return error")
	assert.True(t, exists, "Key should exist")

	// Test Get with non-existent key
	nonExistentKey := "non_existent_key"
	value, err = cacheService.Get(nonExistentKey)
	assert.Error(t, err, "Get should return error for non-existent key")
	assert.Nil(t, value, "Value should be nil for non-existent key")

	// Test Exists with non-existent key
	exists, err = cacheService.Exists(nonExistentKey)
	assert.NoError(t, err, "Exists should not return error for non-existent key")
	assert.False(t, exists, "Non-existent key should not exist")

	// Test Delete method
	err = cacheService.Delete(testKey)
	assert.NoError(t, err, "Delete should not return error")

	// Verify deletion
	exists, err = cacheService.Exists(testKey)
	assert.NoError(t, err, "Exists should not return error after deletion")
	assert.False(t, exists, "Key should not exist after deletion")

	// Test Flush method
	// First add some data
	cacheService.Set("key1", "value1", expiration)
	cacheService.Set("key2", "value2", expiration)

	// Verify data exists
	exists, _ = cacheService.Exists("key1")
	assert.True(t, exists, "key1 should exist before flush")
	exists, _ = cacheService.Exists("key2")
	assert.True(t, exists, "key2 should exist before flush")

	// Flush cache
	err = cacheService.Flush()
	assert.NoError(t, err, "Flush should not return error")

	// Verify all data is gone
	exists, _ = cacheService.Exists("key1")
	assert.False(t, exists, "key1 should not exist after flush")
	exists, _ = cacheService.Exists("key2")
	assert.False(t, exists, "key2 should not exist after flush")
}

// TestServiceFactoryFunctions اختبار دوال إنشاء الخدمات
func TestServiceFactoryFunctions(t *testing.T) {
	// Test all factory functions with nil database
	assert.NotPanics(t, func() {
		authService := NewAuthService(nil)
		assert.NotNil(t, authService, "Auth service should be created")

		userService := NewUserService(nil)
		assert.NotNil(t, userService, "User service should be created")

		serviceService := NewServiceService(nil)
		assert.NotNil(t, serviceService, "Service service should be created")

		categoryService := NewCategoryService(nil)
		assert.NotNil(t, categoryService, "Category service should be created")

		orderService := NewOrderService(nil)
		assert.NotNil(t, orderService, "Order service should be created")

		paymentService := NewPaymentService(nil)
		assert.NotNil(t, paymentService, "Payment service should be created")

		uploadService := NewUploadService(nil)
		assert.NotNil(t, uploadService, "Upload service should be created")

		notificationService := NewNotificationService(nil)
		assert.NotNil(t, notificationService, "Notification service should be created")

		adminService := NewAdminService(nil)
		assert.NotNil(t, adminService, "Admin service should be created")

		cacheService := NewCacheService()
		assert.NotNil(t, cacheService, "Cache service should be created")
	})
}

// TestServiceInterfaces اختبار أن الخدمات تنفذ الواجهات المطلوبة
func TestServiceInterfaces(t *testing.T) {
	container := NewServiceContainer(nil)

	// Test AuthService interface
	var authService AuthService = container.Auth
	assert.NotNil(t, authService, "Auth service should implement AuthService interface")

	// Test UserService interface
	var userService UserService = container.User
	assert.NotNil(t, userService, "User service should implement UserService interface")

	// Test ServiceService interface
	var serviceService ServiceService = container.Service
	assert.NotNil(t, serviceService, "Service service should implement ServiceService interface")

	// Test CategoryService interface
	var categoryService CategoryService = container.Category
	assert.NotNil(t, categoryService, "Category service should implement CategoryService interface")

	// Test OrderService interface
	var orderService OrderService = container.Order
	assert.NotNil(t, orderService, "Order service should implement OrderService interface")

	// Test PaymentService interface
	var paymentService PaymentService = container.Payment
	assert.NotNil(t, paymentService, "Payment service should implement PaymentService interface")

	// Test UploadService interface
	var uploadService UploadService = container.Upload
	assert.NotNil(t, uploadService, "Upload service should implement UploadService interface")

	// Test NotificationService interface
	var notificationService NotificationService = container.Notification
	assert.NotNil(t, notificationService, "Notification service should implement NotificationService interface")

	// Test AdminService interface
	var adminService AdminService = container.Admin
	assert.NotNil(t, adminService, "Admin service should implement AdminService interface")

	// Test CacheService interface
	var cacheService CacheService = container.Cache
	assert.NotNil(t, cacheService, "Cache service should implement CacheService interface")
}

// TestServiceContainerMethods اختبار طرق حاوية الخدمات
func TestServiceContainerMethods(t *testing.T) {
	// Create a mock database connection
	db := &sql.DB{}

	container := NewServiceContainer(db)
	assert.NotNil(t, container, "Service container should be created")

	// Test InitializeDatabase method
	err := container.InitializeDatabase(nil)
	assert.Error(t, err, "InitializeDatabase should return error with nil context")

	// Test Close method
	err = container.Close()
	assert.NoError(t, err, "Close should not return error")
}

// TestCacheExpiration اختبار انتهاء صلاحية التخزين المؤقت
func TestCacheExpiration(t *testing.T) {
	cacheService := NewCacheService()
	assert.NotNil(t, cacheService)

	// Test with very short expiration
	testKey := "expiring_key"
	testValue := "expiring_value"
	shortExpiration := 100 * time.Millisecond

	// Set value with short expiration
	err := cacheService.Set(testKey, testValue, shortExpiration)
	assert.NoError(t, err, "Set should not return error")

	// Immediately get should work
	value, err := cacheService.Get(testKey)
	assert.NoError(t, err, "Get should not return error immediately")
	assert.Equal(t, testValue, value, "Value should match")

	// Wait for expiration
	time.Sleep(200 * time.Millisecond)

	// After expiration, get should fail
	value, err = cacheService.Get(testKey)
	assert.Error(t, err, "Get should return error after expiration")
	assert.Nil(t, value, "Value should be nil after expiration")

	// Exists should return false
	exists, err := cacheService.Exists(testKey)
	assert.NoError(t, err, "Exists should not return error")
	assert.False(t, exists, "Key should not exist after expiration")
}

// TestConcurrentCacheAccess اختبار الوصول المتزامن للتخزين المؤقت
func TestConcurrentCacheAccess(t *testing.T) {
	cacheService := NewCacheService()
	assert.NotNil(t, cacheService)

	// Number of goroutines to run concurrently
	numGoroutines := 10
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			key := "concurrent_key"
			value := "value_from_goroutine_" + string(rune(index))

			// Each goroutine sets a value
			err := cacheService.Set(key, value, 1*time.Minute)
			assert.NoError(t, err, "Set should not return error in goroutine")

			// Each goroutine tries to get the value
			retrieved, err := cacheService.Get(key)
			if err == nil {
				assert.NotNil(t, retrieved, "Retrieved value should not be nil")
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}

// TestServiceErrors اختبار الأخطاء في الخدمات
func TestServiceErrors(t *testing.T) {
	// Test cache service with nil value
	cacheService := NewCacheService()

	// Test Set with empty key
	err := cacheService.Set("", "value", 1*time.Minute)
	assert.NoError(t, err, "Set with empty key should not return error")

	// Test Get with empty key
	value, err := cacheService.Get("")
	assert.Error(t, err, "Get with empty key should return error")
	assert.Nil(t, value, "Value should be nil for empty key")

	// Test Delete with empty key
	err = cacheService.Delete("")
	assert.NoError(t, err, "Delete with empty key should not return error")

	// Test Exists with empty key
	exists, err := cacheService.Exists("")
	assert.NoError(t, err, "Exists with empty key should not return error")
	assert.False(t, exists, "Empty key should not exist")
}

// TestServiceContainerIntegration اختبار تكامل حاوية الخدمات
func TestServiceContainerIntegration(t *testing.T) {
	// Create service container
	container := NewServiceContainer(nil)

	// Verify all services are properly linked
	assert.NotNil(t, container.Auth, "Auth service should be available")
	assert.NotNil(t, container.User, "User service should be available")
	assert.NotNil(t, container.Service, "Service service should be available")
	assert.NotNil(t, container.Category, "Category service should be available")
	assert.NotNil(t, container.Order, "Order service should be available")
	assert.NotNil(t, container.Payment, "Payment service should be available")
	assert.NotNil(t, container.Upload, "Upload service should be available")
	assert.NotNil(t, container.Notification, "Notification service should be available")
	assert.NotNil(t, container.Admin, "Admin service should be available")
	assert.NotNil(t, container.Cache, "Cache service should be available")
}

// TestCreateTablesSQL اختبار إنشاء جداول SQL
func TestCreateTablesSQL(t *testing.T) {
	// Get SQL statements for creating tables
	sqlStatements := CreateTablesSQL()

	// Should return multiple SQL statements
	assert.Greater(t, len(sqlStatements), 0, "Should return at least one SQL statement")

	// Each statement should be a non-empty string
	for i, stmt := range sqlStatements {
		assert.NotEmpty(t, stmt, "SQL statement %d should not be empty", i)
		assert.Contains(t, stmt, "CREATE TABLE", "SQL statement %d should be a CREATE TABLE statement", i)
	}

	// Should include all expected tables
	expectedTables := []string{
		"users",
		"categories",
		"services",
		"orders",
		"payments",
		"notifications",
		"files",
		"system_logs",
	}

	allStatements := ""
	for _, stmt := range sqlStatements {
		allStatements += stmt + " "
	}

	for _, table := range expectedTables {
		assert.Contains(t, allStatements, table, "Should include CREATE TABLE for %s", table)
	}
}

// TestHelperFunctions اختبار الدوال المساعدة
func TestHelperFunctions(t *testing.T) {
	// Test generateID
	id1 := generateID("test")
	id2 := generateID("test")
	assert.NotEmpty(t, id1, "Generated ID should not be empty")
	assert.NotEmpty(t, id2, "Second generated ID should not be empty")
	assert.NotEqual(t, id1, id2, "Generated IDs should be unique")

	// Test serializeStrings and deserializeStrings
	testStrings := []string{"one", "two", "three"}
	serialized := serializeStrings(testStrings)
	assert.NotEmpty(t, serialized, "Serialized string should not be empty")

	deserialized, err := deserializeStrings(serialized)
	assert.NoError(t, err, "Deserialize should not return error")
	assert.Equal(t, testStrings, deserialized, "Deserialized strings should match original")

	// Test with empty array
	emptySerialized := serializeStrings([]string{})
	emptyDeserialized, err := deserializeStrings(emptySerialized)
	assert.NoError(t, err, "Deserialize empty should not return error")
	assert.Empty(t, emptyDeserialized, "Deserialized empty array should be empty")

	// Test validatePaginationParams
	page, limit := validatePaginationParams(0, 0)
	assert.Equal(t, 1, page, "Page should default to 1")
	assert.Equal(t, 10, limit, "Limit should default to 10")

	page, limit = validatePaginationParams(5, 50)
	assert.Equal(t, 5, page, "Page should be 5")
	assert.Equal(t, 50, limit, "Limit should be 50")

	page, limit = validatePaginationParams(-1, 200)
	assert.Equal(t, 1, page, "Negative page should default to 1")
	assert.Equal(t, 100, limit, "Limit above 100 should be capped at 100")

	// Test calculateOffset
	offset := calculateOffset(1, 10)
	assert.Equal(t, 0, offset, "Offset for page 1 should be 0")

	offset = calculateOffset(2, 10)
	assert.Equal(t, 10, offset, "Offset for page 2 should be 10")

	offset = calculateOffset(3, 20)
	assert.Equal(t, 40, offset, "Offset for page 3 with limit 20 should be 40")
}