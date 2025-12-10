package db

import (
	"context"
	"fmt"

	"github.com/nawthtech/nawthtech/backend/internal/config"
)

// إزالة MongoDB والتركيز على Worker API فقط
var (
	// يمكن إزالة هذا الملف بالكامل إذا لم تعد تحتاجه
)

func InitializeSQL(cfg *config.Config, driver, dsn string) error {
	// إزالة تهيئة قاعدة البيانات المحلية
	// سنستخدم Worker API بدلاً من ذلك
	fmt.Println("⚠️ Local database initialization disabled. Using Worker API.")
	return nil
}

func Close() {
	// لا يوجد اتصال قاعدة بيانات للإغلاق
}

// هذه الوظائف لن تعمل بعد الآن
func GetDB() interface{} {
	return nil
}