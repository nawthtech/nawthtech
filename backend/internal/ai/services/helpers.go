package services

import (
    "context"
    "strings"
)

// UserContextKey مفتاح سياق المستخدم
type UserContextKey string

const (
    UserIDKey   UserContextKey = "user_id"
    UserTierKey UserContextKey = "user_tier"
    UserEmailKey UserContextKey = "user_email"
    UserRoleKey  UserContextKey = "user_role"
)

// extractUserIDFromContext استخراج معرف المستخدم من السياق
func extractUserIDFromContext(ctx context.Context) string {
    if ctx == nil {
        return ""
    }
    
    // المحاولة الأولى: باستخدام المفاتيح المخصصة
    if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
        return userID
    }
    
    // المحاولة الثانية: باستخدام مفتاح عام
    if userID, ok := ctx.Value("userID").(string); ok && userID != "" {
        return userID
    }
    
    // المحاولة الثالثة: باستخدام مفتاح UserId
    if userID, ok := ctx.Value("UserId").(string); ok && userID != "" {
        return userID
    }
    
    return ""
}

// extractUserTierFromContext استخراج طبقة المستخدم من السياق
func extractUserTierFromContext(ctx context.Context) string {
    if ctx == nil {
        return "free"
    }
    
    // المحاولة الأولى: باستخدام المفاتيح المخصصة
    if userTier, ok := ctx.Value(UserTierKey).(string); ok && userTier != "" {
        return strings.ToLower(userTier)
    }
    
    // المحاولة الثانية: باستخدام مفتاح عام
    if userTier, ok := ctx.Value("userTier").(string); ok && userTier != "" {
        return strings.ToLower(userTier)
    }
    
    // المحاولة الثالثة: محاولة استنتاج الطبقة من الدور
    if userRole, ok := ctx.Value(UserRoleKey).(string); ok && userRole != "" {
        switch strings.ToLower(userRole) {
        case "premium", "vip", "gold":
            return "premium"
        case "basic", "silver", "standard":
            return "basic"
        case "admin", "administrator":
            return "premium" // المسؤولون يحصلون على مميزات ممتازة
        default:
            return "free"
        }
    }
    
    return "free"
}

// extractUserEmailFromContext استخراج البريد الإلكتروني للمستخدم من السياق
func extractUserEmailFromContext(ctx context.Context) string {
    if ctx == nil {
        return ""
    }
    
    if email, ok := ctx.Value(UserEmailKey).(string); ok && email != "" {
        return email
    }
    
    if email, ok := ctx.Value("userEmail").(string); ok && email != "" {
        return email
    }
    
    return ""
}

// extractUserRoleFromContext استخراج دور المستخدم من السياق
func extractUserRoleFromContext(ctx context.Context) string {
    if ctx == nil {
        return "user"
    }
    
    if role, ok := ctx.Value(UserRoleKey).(string); ok && role != "" {
        return role
    }
    
    if role, ok := ctx.Value("userRole").(string); ok && role != "" {
        return role
    }
    
    return "user"
}

// WithUserContext إضافة معلومات المستخدم إلى السياق
func WithUserContext(ctx context.Context, userID, userTier, email, role string) context.Context {
    if ctx == nil {
        ctx = context.Background()
    }
    
    ctx = context.WithValue(ctx, UserIDKey, userID)
    ctx = context.WithValue(ctx, UserTierKey, userTier)
    ctx = context.WithValue(ctx, UserEmailKey, email)
    ctx = context.WithValue(ctx, UserRoleKey, role)
    
    return ctx
}

// IsAuthenticated التحقق إذا كان المستخدم مصادقاً
func IsAuthenticated(ctx context.Context) bool {
    return extractUserIDFromContext(ctx) != ""
}

// IsPremiumUser التحقق إذا كان المستخدم من الطبقة الممتازة
func IsPremiumUser(ctx context.Context) bool {
    tier := extractUserTierFromContext(ctx)
    return tier == "premium" || tier == "admin"
}

// IsAdminUser التحقق إذا كان المستخدم مسؤولاً
func IsAdminUser(ctx context.Context) bool {
    role := extractUserRoleFromContext(ctx)
    return strings.ToLower(role) == "admin" || strings.ToLower(role) == "administrator"
}

// GetUserContextInfo الحصول على جميع معلومات المستخدم من السياق
func GetUserContextInfo(ctx context.Context) map[string]string {
    return map[string]string{
        "user_id":    extractUserIDFromContext(ctx),
        "user_tier":  extractUserTierFromContext(ctx),
        "user_email": extractUserEmailFromContext(ctx),
        "user_role":  extractUserRoleFromContext(ctx),
        "authenticated": toString(IsAuthenticated(ctx)),
        "premium":      toString(IsPremiumUser(ctx)),
        "admin":        toString(IsAdminUser(ctx)),
    }
}

// toString تحويل قيمة منطقية إلى سلسلة نصية
func toString(b bool) string {
    if b {
        return "true"
    }
    return "false"
}

// ValidateUserContext التحقق من صحة سياق المستخدم
func ValidateUserContext(ctx context.Context) error {
    userID := extractUserIDFromContext(ctx)
    if userID == "" {
        return ErrUnauthorized
    }
    
    // يمكن إضافة مزيد من التحقق هنا
    // مثل: التحقق من صلاحية التوكن، تاريخ انتهاء الصلاحية، إلخ.
    
    return nil
}

// ErrUnauthorized خطأ المصادقة
var ErrUnauthorized = &ServiceError{"unauthorized", "User is not authenticated"}